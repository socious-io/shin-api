package views

import (
	"context"
	"fmt"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/database"
	"shin/src/lib"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func createDiscordReviewMessage(kyb *models.KYBVerification, u *models.User, org *models.Organization) string {

	documents := ""
	for i, document := range kyb.Documents {
		documents = fmt.Sprintf("%s\n%v. %s", documents, i, document.Url)
	}

	message := fmt.Sprintf("ID: %s\n", kyb.ID)
	message += "\nUser--------------------------------\n"
	message += fmt.Sprintf("ID: %s\n", u.ID)
	message += fmt.Sprintf("Firstname: %s\n", *u.FirstName)
	message += fmt.Sprintf("Lastname: %s\n", *u.LastName)
	message += fmt.Sprintf("Email: %s\n", u.Email)
	message += "\nOrganization------------------------\n"
	message += fmt.Sprintf("ID: %s\n", org.ID)
	message += fmt.Sprintf("Name: %s\n", org.Name)
	message += fmt.Sprintf("Description: %s\n", org.Description)
	message += fmt.Sprintf("\nDocuments---------------------------%s\n\n", documents)
	message += fmt.Sprintf("\nReviewing----------------------------\n")
	message += fmt.Sprintf("Approve: <%s/kyb/%s/approve?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.Admin.AccessToken)
	message += fmt.Sprintf("Reject: <%s/kyb/%s/reject?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.Admin.AccessToken)

	return message

}

func kybVerificationGroup(router *gin.Engine) {
	g := router.Group("kyb")

	g.POST("/:org_id", auth.LoginRequired(), func(c *gin.Context) {
		orgID := c.Param("org_id")
		u, _ := c.Get("user")
		ctx, _ := c.Get("ctx")

		form := new(KYBVerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, err := models.GetOrgByMember(uuid.MustParse(orgID), u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Organization not found",
			})
			return
		}

		kyb := &models.KYBVerification{
			UserID: u.(*models.User).ID,
			OrgID:  org.ID,
		}

		kyb, err = kyb.Create(ctx.(context.Context), form.Documents)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lib.DiscordSendTextMessage(
			config.Config.Logger.Discord["shin_kyb_channel"],
			createDiscordReviewMessage(kyb, u.(*models.User), org),
		)

		c.JSON(http.StatusOK, kyb)
	})

	g.GET("/", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		u, _ := c.Get("user")
		paginate, _ := c.Get("paginate")
		limit, _ := c.Get("limit")
		page, _ := c.Get("page")

		kybVerifications, total, err := models.GetAllByUserId(u.(*models.User).ID, paginate.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": kybVerifications,
			"page":    page,
			"limit":   limit,
			"total":   total,
		})
	})

	g.GET("/:id", auth.LoginRequired(), func(c *gin.Context) {
		u, _ := c.Get("user")
		verificationId := c.Param("id")

		verification, err := models.GetByIdAndUserId(uuid.MustParse(verificationId), u.(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, verification)
	})

	g.GET("/:id/approve", adminAccessRequired(), func(c *gin.Context) {

		ctx, _ := c.Get("ctx")
		verificationId := c.Param("id")

		verification, err := models.GetById(uuid.MustParse(verificationId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := verification.ChangeStatus(ctx.(context.Context), models.KYBStatusApproved); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, err := models.GetOrg(verification.OrgID)

		if err := org.UpdateVerification(ctx.(context.Context), true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	g.GET("/:id/reject", adminAccessRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		verificationId := c.Param("id")

		kybVerification := models.KYBVerification{
			ID: uuid.MustParse(verificationId),
		}

		if err := kybVerification.ChangeStatus(ctx.(context.Context), models.KYBStatusRejected); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

}
