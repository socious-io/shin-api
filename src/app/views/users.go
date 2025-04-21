package views

import (
	"net/http"
	"shin/src/app/auth"

	"github.com/gin-gonic/gin"
)

func userGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("", func(c *gin.Context) {
		u, _ := c.Get("user")
		c.JSON(http.StatusOK, u)
	})

}
