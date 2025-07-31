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
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		schemas, total, err := models.GetSchemas(u.(*models.User).ID, page.(database.Paginate))
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
		// u, _ := c.Get("user")
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
		u, _ := c.Get("user")
		s.CreatedID = &u.(*models.User).ID
		ctx, _ := c.Get("ctx")
		if err := s.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		u, _ := c.Get("user")
		s, err := models.GetSchema(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if s.Created.ID != u.(*models.User).ID || !s.Deleteable {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := s.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
