package views

import (
	"context"
	"net/http"
	"net/url"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/utils"
	"strings"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func verificationsGroup(router *gin.Engine) {
	g := router.Group("verifications")

	g.GET("", paginate(), LoginRequired(), func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		page := c.MustGet("paginate").(database.Paginate)
		verifications, total, err := models.GetVerifications(u.ID, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": verifications,
			"total":   total,
		})
	})

	g.GET("/:id/individuals", paginate(), LoginRequired(), func(c *gin.Context) {
		u := c.MustGet("user").(*models.User)
		page := c.MustGet("paginate").(database.Paginate)
		id := c.Param("id")
		verifications, total, err := models.GetVerificationsIndividuals(u.ID, uuid.MustParse(id), page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": verifications,
			"total":   total,
		})
	})

	g.GET("/:id", LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.GET("/individuals/:id", LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerificationsIndividual(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.POST("individuals", AuthRequired(), func(c *gin.Context) {
		form := new(VerificationIndividualForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := new(models.VerificationIndividual)
		u := c.MustGet("user").(*models.User)
		v.UserID = u.ID
		v.VerificationID = form.VerificationID
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := v.Create(ctx, form.CustomerID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.GET("/:id/connect", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerificationsIndividual(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if v.ConnectionURL != nil {
			if time.Since(*v.ConnectionAt) < 2*time.Minute {
				c.JSON(http.StatusOK, v)
				return
			}
		}
		ctx, _ := c.MustGet("ctx").(context.Context)

		callback, _ := url.JoinPath(config.Config.Host, strings.ReplaceAll(c.Request.URL.String(), "connect", "callback"))

		if err := v.NewConnection(ctx, callback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, v)
	})

	g.GET("/:id/callback", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerificationsIndividual(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := v.ProofRequest(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.GET("/:id/verify", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerificationsIndividual(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		v.ProofVerify(ctx)
		c.JSON(http.StatusOK, v)
	})

	g.POST("", LoginRequired(), func(c *gin.Context) {
		form := new(VerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := new(models.Verification)
		utils.Copy(form, v)
		u := c.MustGet("user").(*models.User)
		v.UserID = u.ID
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := v.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.PUT("/:id", LoginRequired(), func(c *gin.Context) {
		form := new(VerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := c.MustGet("user").(*models.User)
		if v.UserID.String() != u.ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		utils.Copy(form, v)

		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := v.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, v)
	})

	g.DELETE("/:id", LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := c.MustGet("user").(*models.User)
		if v.UserID.String() != u.ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := v.Delete(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
