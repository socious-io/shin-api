package views

import (
	"context"
	"net/http"
	"net/url"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/database"
	"shin/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func verificationsGroup(router *gin.Engine) {
	g := router.Group("verifications")

	g.GET("", paginate(), auth.LoginRequired(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		verifications, total, err := models.GetVerifications(u.(*models.User).ID, page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": verifications,
			"total":   total,
		})
	})

	g.GET("/:id/individuals", paginate(), auth.LoginRequired(), func(c *gin.Context) {
		u, _ := c.Get("user")
		page, _ := c.Get("paginate")
		id := c.Param("id")
		verifications, total, err := models.GetVerificationsIndividuals(u.(*models.User).ID, uuid.MustParse(id), page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": verifications,
			"total":   total,
		})
	})

	g.GET("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.GET("/individuals/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerificationsIndividual(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.POST("individuals", auth.AuthRequired(), func(c *gin.Context) {
		form := new(VerificationIndividualForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := new(models.VerificationIndividual)
		u, _ := c.Get("user")
		v.UserID = u.(*models.User).ID
		v.VerificationID = form.VerificationID
		ctx, _ := c.Get("ctx")
		if err := v.Create(ctx.(context.Context), form.CustomerID); err != nil {
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
		ctx, _ := c.Get("ctx")

		callback, _ := url.JoinPath(config.Config.Host, strings.ReplaceAll(c.Request.URL.String(), "connect", "callback"))

		if err := v.NewConnection(ctx.(context.Context), callback); err != nil {
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
		ctx, _ := c.Get("ctx")
		if err := v.ProofRequest(ctx.(context.Context)); err != nil {
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
		ctx, _ := c.Get("ctx")
		v.ProofVerify(ctx.(context.Context))
		c.JSON(http.StatusOK, v)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		form := new(VerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v := new(models.Verification)
		utils.Copy(form, v)
		u, _ := c.Get("user")
		v.UserID = u.(*models.User).ID
		ctx, _ := c.Get("ctx")
		if err := v.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.PUT("/:id", auth.LoginRequired(), func(c *gin.Context) {
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
		u, _ := c.Get("user")
		if v.UserID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		utils.Copy(form, v)

		ctx, _ := c.Get("ctx")
		if err := v.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, v)
	})

	g.DELETE("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if v.UserID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := v.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
