package views

import (
	"context"
	"net/http"
	"shin/src/app/models"
	"shin/src/utils"

	"github.com/gin-gonic/gin"
)

func userGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(LoginRequired())

	g.GET("", func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		c.JSON(http.StatusOK, u)
	})

	g.PUT("/profile", func(c *gin.Context) {
		form := new(ProfileUpdateForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := c.MustGet("user").(*models.User)
		ctx, _ := c.MustGet("ctx").(context.Context)
		utils.Copy(form, user)

		if err := user.UpdateProfile(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, user)
	})

}
