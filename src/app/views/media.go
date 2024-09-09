package views

import (
	"context"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func mediaGroup(router *gin.Engine) {
	g := router.Group("media")

	g.GET("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		m, err := models.GetMedia(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't get media",
			})
			return
		}
		c.JSON(http.StatusOK, m)
	})

	g.POST("/upload", auth.LoginRequired(), func(c *gin.Context) {

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't upload media",
			})
			return
		}
		defer file.Close()

		mediaUrl, fileName := services.Upload(file, header.Filename)

		u, _ := c.Get("user")

		media := models.Media{
			UserID:   u.(*models.User).ID,
			URL:      mediaUrl,
			Filename: fileName,
		}

		ctx, _ := c.Get("ctx")
		if err := media.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't upload media",
			})
			return
		}

		c.JSON(http.StatusCreated, media)

	})
}