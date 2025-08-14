package views

import (
	"bytes"
	"io"
	"net/http"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/lib"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/socious-io/goaccount"
	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
* Authorization
 */
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		claims, err := goaccount.ClaimsFromBearerToken(tokenStr)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		u, err := models.GetUser(uuid.MustParse(claims.ID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)
		/*
			identityStr := c.GetHeader(http.CanonicalHeaderKey("current-identity"))
			if identityUUID, err := uuid.Parse(identityStr); err == nil {
				c.Set("identity", identityUUID)
			} */
		c.Next()
	}
}
func IntegrationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("apikey")
		if tokenStr == "" {
			tokenStr = c.Query("apikey")
		}
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "apikey is required"})
			c.Abort()
			return
		}
		key, err := models.GetIntegrationBySecret(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid apikey"})
			c.Abort()
			return
		}
		u, err := models.GetUser(key.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		splited := strings.Split(tokenStr, " ")
		if len(splited) > 1 {
			tokenStr = splited[1]
		} else {
			tokenStr = splited[0]
		}
		if tokenStr != "" {
			claims, err := goaccount.VerifyToken(tokenStr)
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
					c.Abort()
					return
				}
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			u, err := models.GetUser(uuid.MustParse(claims.ID))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.Set("user", u)
			c.Next()
			return
		}

		tokenStr = c.GetHeader("apikey")
		if tokenStr == "" {
			tokenStr = c.Query("apikey")
		}
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "apikey is required"})
			c.Abort()
			return
		}
		key, err := models.GetIntegrationBySecret(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid apikey"})
			c.Abort()
			return
		}
		u, err := models.GetUser(key.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)
		c.Next()
	}
}

// Pagination
func paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		if page < 1 {
			page = 1
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}
		filters := make([]database.Filter, 0)
		for key, values := range c.Request.URL.Query() {
			if strings.Contains(key, "filter.") && len(values) > 0 {
				filters = append(filters, database.Filter{
					Key:   strings.Replace(key, "filter.", "", -1),
					Value: values[0],
				})
			}
		}

		c.Set("paginate", database.Paginate{
			Limit:   limit,
			Offet:   (page - 1) * limit,
			Filters: filters,
		})
		c.Set("limit", limit)
		c.Set("page", page)
		c.Next()

	}
}

// Logger
type GinLoggerResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *GinLoggerResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLoggerMiddleware(logger *lib.GinLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &GinLoggerResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		c.Writer = w
		start := time.Now()
		requestId := uuid.NewString()

		// Process request
		c.Next()

		logger.Auto(requestId, lib.GinLogFields{
			Duration:       time.Since(start),
			StatusCode:     w.Status(),
			RequestHeaders: c.Request.Header,
			Headers:        w.Header(),
			RequestBody:    bytes.NewBuffer(requestBody),
			Body:           w.body,
			IP:             c.ClientIP(),
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			Query:          c.Request.URL.RawQuery,
		})
	}
}

// Administration
func AdminAccessRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		access_token := c.Query("admin_access_token")
		isAdmin := access_token == config.Config.Admin.AccessToken

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "AdminAccessRequired"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AccountCenterRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("x-account-center-id")
		secret := c.Request.Header.Get("x-account-center-secret")
		hash, _ := goaccount.HashPassword(secret)

		if id != config.Config.GoAccounts.ID || goaccount.CheckPasswordHash(secret, hash) != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account center required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
