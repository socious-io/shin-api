package views

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func syncGroup(router *gin.Engine) {
	g := router.Group("sync")
	g.Use(AccountCenterRequired())

	g.PUT("", func(c *gin.Context) {
		form := new(SyncForm)

		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx := c.MustGet("ctx").(context.Context)
		if err := form.User.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, o := range form.Organizations {
			if err := o.Create(ctx, form.User.ID); err != nil {
				log.Println(err.Error(), o)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}
