package views

import (
	"bytes"
	"io"
	"net/http"
	"shin/src/config"
	"shin/src/database"
	"shin/src/lib"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
func adminAccessRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		access_token := c.Query("admin_access_token")
		isAdmin := access_token == config.Config.Admin.AccessToken

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "AdminAccessRequired"})
			c.Abort()
			return
		}

		c.Next()
		return
	}
}
