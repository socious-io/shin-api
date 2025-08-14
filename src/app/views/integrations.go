package views

import (
	"context"
	"net/http"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/lib"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func integrationGroup(router *gin.Engine) {
	g := router.Group("integrations")
	g.Use(LoginRequired())

	g.GET("/keys", paginate(), LoginRequired(), func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		paginate, _ := c.MustGet("paginate").(database.Paginate)
		limit, _ := c.Get("limit")
		page, _ := c.Get("page")

		integrationKeys, total, err := models.GetIntegrations(u.ID, paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": integrationKeys,
			"page":    page,
			"limit":   limit,
			"total":   total,
		})
	})

	g.POST("/keys", func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		ctx, _ := c.MustGet("ctx").(context.Context)

		form := new(models.IntegrationKey)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		integrationKey := &models.IntegrationKey{
			UserID:  u.ID,
			Name:    form.Name,
			BaseUrl: config.Config.Host,
			Key:     lib.GenerateApiKey(),
			Secret:  lib.GenerateApiSecret(),
		}

		integrationKeyCreated, err := integrationKey.Create(ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, integrationKeyCreated)
	})

	g.PUT("/keys/:id", func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := c.Param("id")

		integrationKey := new(models.IntegrationKey)
		if err := c.ShouldBindJSON(integrationKey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		integrationKey.ID = uuid.MustParse(id)

		if err := integrationKey.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, integrationKey)
	})

	g.DELETE("/keys/:id", func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := c.Param("id")

		integrationKey := &models.IntegrationKey{
			ID:     uuid.MustParse(id),
			UserID: u.ID,
		}

		if err := integrationKey.Delete(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
