package views

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/database"
	"shin/src/lib"
	"shin/src/services"
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

	g.PATCH("/revoke", auth.LoginRequired(), func(c *gin.Context) {

		u, _ := c.Get("user")

		form := new(CredentialBulkOperationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		credentials, err := models.GetCredentialsByIds(form.Credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(credentials) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "there's no matching credential(s)"})
			return
		}

		//Validate credential(s) ownerships
		for _, credential := range credentials {
			if credential.CreatedID.String() != u.(*models.User).ID.String() {
				c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
				return
			}
		}

		//Handling revoke async
		for _, credential := range credentials {
			go services.SendOperation(services.OperationConfig{
				Trigger: models.OperationCredentialRevoke,
				Entity:  credential,
			})
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

		//Sending Email
		if cv.Recipient != nil && cv.Recipient.Email != nil {
			items := map[string]string{
				"title":      cv.Name,
				"issuer_org": cv.Organization.Name,
				"recipient":  fmt.Sprintf("%s %s", *cv.Recipient.FirstName, *cv.Recipient.LastName),
				"link":       fmt.Sprintf("%s/connect/credential/%s", config.Config.FrontHost, cv.ID.String()),
			}
			services.SendEmail(services.EmailConfig{
				Approach:    services.EmailApproachTemplate,
				Destination: *cv.Recipient.Email,
				Title:       "Shin - Your verification credentials",
				Template:    "credentials-recipients",
				Args:        items,
			})
		}

		c.JSON(http.StatusCreated, cv)
	})

	g.POST("/with-recipient", auth.LoginRequired(), func(c *gin.Context) {

		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")

		form := new(CredentialRecipientForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Creating Recipient
		r := new(models.Recipient)
		utils.Copy(form.Recipient, r)
		r.UserID = u.(*models.User).ID
		if err := r.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Creating Credential
		schema, err := models.GetSchema(form.Credential.SchemaID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if schema.IssueDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "schema for issuing credentials is disabled"})
			return
		}
		cv := new(models.Credential)
		cv.CreatedID = u.(*models.User).ID
		orgs, err := models.GetOrgsByMember(cv.CreatedID)
		if err != nil || len(orgs) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("fetching org error :%v", err)})
			return
		}
		utils.Copy(form.Credential, cv)
		cv.OrganizationID = orgs[0].ID
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

	g.POST("/import", auth.LoginRequired(), func(c *gin.Context) {

		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")
		user := u.(*models.User)

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't upload media",
			})
			return
		}
		defer file.Close()

		SchemaID := c.Request.FormValue("schema_id")

		// Fetching Schema
		schema, err := models.GetSchema(uuid.MustParse(SchemaID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if schema.IssueDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "schema for issuing credentials is disabled"})
			return
		}

		schemaAttributes := map[string]string{}
		for _, attributes := range schema.Attributes {
			schemaAttributes[attributes.Name] = string(attributes.Type)
		}
		schemaAttributes["recipient_first_name"] = string(models.Text)
		schemaAttributes["recipient_last_name"] = string(models.Text)
		schemaAttributes["recipient_email"] = string(models.Email)

		//Processing CSV file
		resultChan, errChan := make(chan []map[string]any), make(chan error)

		go lib.ValidateCSVFile(file, schemaAttributes, resultChan, errChan)

		for {
			select {
			case results := <-resultChan:
				i := models.Import{
					Target:     models.ImportTargetCredentials,
					UserID:     user.ID,
					TotalCount: len(results),
				}
				if err := i.Create(ctx.(context.Context)); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, i)
				go services.InitiateImport(results, map[string]any{
					"schema_id": schema.ID,
					"user_id":   user.ID,
					"file_name": header.Filename,
				}, i)
				return
			case err := <-errChan:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

	})

	g.GET("/import/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")

		i, err := models.GetImport(uuid.MustParse(id))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, i)
	})

	g.GET("/import/download-sample/:schema_id", func(c *gin.Context) {

		SchemaID := c.Param("schema_id")

		// Fetching Schema
		schema, err := models.GetSchema(uuid.MustParse(SchemaID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if schema.IssueDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "schema for issuing credentials is disabled"})
			return
		}

		schemaAttributes := []string{"recipient_first_name", "recipient_last_name", "recipient_email"}
		schemaFields := []string{"Recipient First name", "Recipient Last name", "recipient@email.com"}
		for _, attributes := range schema.Attributes {
			attribute := attributes.Name
			attribute_type := string(attributes.Type)
			var sample_value string

			switch attribute_type {
			case string(models.Text):
				sample_value = "some text"
				break
			case string(models.Number):
				sample_value = "1234"
				break
			case string(models.Boolean):
				sample_value = "true"
				break
			case string(models.Email):
				sample_value = "example@email.com"
				break
			case string(models.Url):
				sample_value = "http://some.url.example"
				break
			case string(models.Datetime):
				sample_value = string(time.Now().Format(time.RFC3339))
				break
			default:
				sample_value = "UNKNOWN_DATATYPE"
				break
			}

			schemaAttributes = append(schemaAttributes, attribute)
			schemaFields = append(schemaFields, sample_value)
		}

		// Set headers for CSV download
		c.Header("Content-Disposition", "attachment;filename=sample-import.csv")
		c.Header("Content-Type", "text/csv")

		// Create a CSV writer that writes to the response writer
		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()
		if err := writer.Write(schemaAttributes); err != nil {
			c.String(500, "Could not write CSV header: %v", err)
			return
		}
		if err := writer.Write(schemaFields); err != nil {
			c.String(500, "Could not write CSV field: %v", err)
			return
		}

	})

	g.POST("/notify", auth.LoginRequired(), func(c *gin.Context) {
		form := new(CredentialBulkEmailForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		credentials, err := models.GetCredentialsByIds(form.Credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Sending Email
		for _, credential := range credentials {
			if credential.Recipient != nil && credential.Recipient.Email != nil {
				items := map[string]string{
					"title":      credential.Name,
					"issuer_org": credential.Organization.Name,
					"recipient":  fmt.Sprintf("%s %s", *credential.Recipient.FirstName, *credential.Recipient.LastName),
					"link":       fmt.Sprintf("%s/connect/credential/%s", config.Config.FrontHost, credential.ID.String()),
					"message":    form.Message,
				}
				services.SendEmail(services.EmailConfig{
					Approach:    services.EmailApproachTemplate,
					Destination: *credential.Recipient.Email,
					Title:       "Shin: Your Verifiable Credential is Here",
					Template:    "credentials-recipients",
					Args:        items,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.POST("/notify/via-schema", auth.LoginRequired(), func(c *gin.Context) {

		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")

		form := new(CredentialBySchemaEmailForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doneChan, errChan := make(chan bool), make(chan error)

		go services.CredentialBulkEmailAsync(ctx.(context.Context), u.(*models.User).ID, form.SchemaID.String(), form.Message, doneChan, errChan)
		for {
			select {
			case <-doneChan:
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
				})
				return
			case err := <-errChan:
				c.JSON(http.StatusOK, gin.H{
					"error": err,
				})
			}
		}

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

	g.POST("/delete", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")

		form := new(CredentialBulkOperationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := models.CredentialsBulkDelete(ctx.(context.Context), form.Credentials, u.(*models.User).ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
