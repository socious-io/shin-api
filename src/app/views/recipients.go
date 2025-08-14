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

func recipientsGroup(router *gin.Engine) {
	g := router.Group("recipients")
	g.Use(AuthRequired())

	g.GET("", paginate(), func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		page := c.MustGet("paginate").(database.Paginate)
		recipients, total, err := models.SearchRecipients(c.Query("q"), u.ID, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": recipients,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		r, err := models.GetRecipient(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, r)
	})

	g.POST("", func(c *gin.Context) {
		form := new(RecipientForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r := new(models.Recipient)
		utils.Copy(form, r)
		u := c.MustGet("user").(*models.User)
		r.UserID = u.ID
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := r.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, r)
	})

	g.PUT("/:id", func(c *gin.Context) {
		form := new(RecipientForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		r, err := models.GetRecipient(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := c.MustGet("user").(*models.User)
		if r.UserID.String() != u.ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		utils.Copy(form, r)

		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := r.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, r)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		r, err := models.GetRecipient(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := c.MustGet("user").(*models.User)
		if r.UserID.String() != u.ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := r.Delete(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
