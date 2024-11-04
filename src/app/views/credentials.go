package views

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/database"
	"shin/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func credentialsGroup(router *gin.Engine) {
	g := router.Group("credentials")

	g.GET("", paginate(), auth.LoginRequired(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		credentials, total, err := models.GetCredentials(u.(*models.User).ID, page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": credentials,
			"total":   total,
		})
	})

	g.GET("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.GET("/:id/connect", func(c *gin.Context) {
		id := c.Param("id")
		cv, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if cv.ConnectionURL != nil {
			if time.Since(*cv.ConnectionAt) < 2*time.Minute {
				c.JSON(http.StatusOK, cv)
				return
			}
		}
		ctx, _ := c.Get("ctx")

		callback, _ := url.JoinPath(config.Config.Host, strings.ReplaceAll(c.Request.URL.String(), "connect", "callback"))

		if err := cv.NewConnection(ctx.(context.Context), callback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cv)
	})

	g.GET("/:id/callback", func(c *gin.Context) {
		id := c.Param("id")
		cv, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := cv.Issue(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.PATCH("/:id/revoke", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		cv, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if cv.CreatedID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := cv.Revoke(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		form := new(CredentialForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		schema, err := models.GetSchema(form.SchemaID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if schema.IssueDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "schema for issuing credentials is disabled"})
			return
		}
		cv := new(models.Credential)
		u, _ := c.Get("user")
		cv.CreatedID = u.(*models.User).ID
		ctx, _ := c.Get("ctx")
		orgs, err := models.GetOrgsByMember(cv.CreatedID)
		if err != nil || len(orgs) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("fetching org error :%v", err)})
			return
		}
		utils.Copy(form, cv)
		cv.OrganizationID = orgs[0].ID
		claims := gin.H{}
		for _, claim := range form.Claims {
			claims[claim.Name] = claim.Value
		}
		claims["type"] = schema.Name
		claims["issued_date"] = time.Now().Format(time.RFC3339)
		claims["company_name"] = orgs[0].Name
		cv.Claims, _ = json.Marshal(&claims)
		if err := cv.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, cv)
	})

	g.POST("/with-recipient", auth.LoginRequired(), func(c *gin.Context) {

		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")

		form := new(CredentialRecipientForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		//Creating Recipient
		r := new(models.Recipient)
		utils.Copy(form.Recipient, r)
		r.UserID = u.(*models.User).ID
		if err := r.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}

		//Creating Credential
		schema, err := models.GetSchema(form.Credential.SchemaID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
			return
		}
		if schema.IssueDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error4": "schema for issuing credentials is disabled"})
			return
		}
		cv := new(models.Credential)
		cv.CreatedID = u.(*models.User).ID
		orgs, err := models.GetOrgsByMember(cv.CreatedID)
		if err != nil || len(orgs) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error5": fmt.Sprintf("fetching org error :%v", err)})
			return
		}
		utils.Copy(form.Credential, cv)
		cv.OrganizationID = orgs[0].ID
		fmt.Println(r.ID)
		cv.RecipientID = &r.ID
		claims := gin.H{}
		for _, claim := range form.Credential.Claims {
			claims[claim.Name] = claim.Value
		}
		claims["type"] = schema.Name
		claims["issued_date"] = time.Now().Format(time.RFC3339)
		claims["company_name"] = orgs[0].Name
		cv.Claims, _ = json.Marshal(&claims)
		if err := cv.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, cv)
	})

	g.PUT("/:id", auth.LoginRequired(), func(c *gin.Context) {
		form := new(CredentialForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := c.Param("id")
		cv, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if cv.CreatedID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}

		if cv.Status == models.StatusClaimed {
			c.JSON(http.StatusForbidden, gin.H{"error": "no update allowed after claim"})
			return
		}
		utils.Copy(form, cv)

		ctx, _ := c.Get("ctx")
		if err := cv.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, cv)
	})

	g.DELETE("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		cv, err := models.GetCredential(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if cv.CreatedID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := cv.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
