package views

import (
	"context"
	"net/http"
	"shin/src/app/models"
	"shin/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func orgGroup(router *gin.Engine) {
	g := router.Group("organizations")
	g.Use(LoginRequired())

	g.GET("", func(c *gin.Context) {
		u, _ := c.Get("user")

		orgs, err := models.GetOrgsByMember(u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"results": orgs})
	})

	g.GET("/:id", func(c *gin.Context) {
		orgID := c.Param("id")
		org, err := models.GetOrg(uuid.MustParse(orgID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, org)
	})

	g.POST("", func(c *gin.Context) {
		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		o := new(models.Organization)
		utils.Copy(form, o)
		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")
		if err := o.Create(ctx.(context.Context), u.(*models.User).ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, o)
	})

	g.PUT("/:id", func(c *gin.Context) {
		orgID := c.Param("id")
		u, _ := c.Get("user")
		// TODO: can be middleware
		o, err := models.GetOrgByMember(uuid.MustParse(orgID), u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.Copy(form, o)
		ctx, _ := c.Get("ctx")
		if err := o.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, o)
	})
}
