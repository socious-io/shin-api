package views

import (
	"context"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/utils"

	"github.com/gin-gonic/gin"
)

func userGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("", func(c *gin.Context) {
		u, _ := c.Get("user")
		c.JSON(http.StatusOK, u)
	})

	g.PUT("/profile", func(c *gin.Context) {
		form := new(ProfileUpdateForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")
		user := u.(*models.User)
		utils.Copy(form, user)

		if err := user.UpdateProfile(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, u)
	})

}
