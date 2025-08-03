package views

import (
	"context"
	"net/http"
	"shin/src/app/models"
	"shin/src/utils"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func schemasGroup(router *gin.Engine) {
	g := router.Group("schemas")
	g.Use(LoginRequired())

	g.GET("", paginate(), func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		page := c.MustGet("paginate").(database.Paginate)
		schemas, total, err := models.GetSchemas(u.ID, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{
			"results": schemas,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		// u := c.MustGet("user").(*models.User)
		s, err := models.GetSchema(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, s)
	})

	g.POST("", func(c *gin.Context) {
		form := new(SchemaForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s := new(models.Schema)
		utils.Copy(form, s)
		u := c.MustGet("user").(*models.User)
		s.CreatedID = &u.ID
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := s.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		u := c.MustGet("user").(*models.User)
		s, err := models.GetSchema(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if s.Created.ID != u.ID || !s.Deleteable {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := s.Delete(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
