package views

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"regexp"
	"shin/src/app/models"
	"shin/src/services"
	"shin/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/socious-io/goaccount"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	//Socious ID
	g.POST("", func(c *gin.Context) {
		form := new(AuthForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionPolicies := []goaccount.PolicyType{
			goaccount.PolicyTypePreventUserAccountSelection,
			goaccount.PolicyTypeRequireAtleastOneOrg,
		}
		session, authURL, err := goaccount.StartSession(
			form.RedirectURL,
			goaccount.AuthModeLogin,
			sessionPolicies,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"session":  session,
			"auth_url": authURL,
		})
	})

	g.POST("/session", func(c *gin.Context) {
		form := new(SessionForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := goaccount.GetSessionToken(form.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var (
			connect *models.OauthConnect
			user    = new(models.User)
			ctx     = c.MustGet("ctx").(context.Context)
		)
		u, err := token.GetUserProfile()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.Copy(u, user)

		if connect, err = models.GetOauthConnectByMUI(user.ID.String(), models.OauthConnectedProvidersSociousID); err != nil {
			connect = &models.OauthConnect{
				Provider:       models.OauthConnectedProvidersSociousID,
				AccessToken:    token.AccessToken,
				RefreshToken:   &token.RefreshToken,
				MatrixUniqueID: user.ID.String(),
				UserId:         user.ID,
			}
		}

		if err := user.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgs, err := token.GetMyOrganizations()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, org := range orgs {
			var o = new(models.Organization)
			utils.Copy(org, o)
			if err := o.Create(ctx, user.ID); err != nil {
				log.Println(err.Error(), o)
			}
		}

		if err := connect.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := goaccount.GenerateFullTokens(user.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, jwt)
	})

	g.POST("/login", func(c *gin.Context) {
		form := new(LoginForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}
		if u.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}
		if err := goaccount.CheckPasswordHash(form.Password, *u.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}

		tokens, err := goaccount.GenerateFullTokens(u.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
	})

	g.POST("/register", func(c *gin.Context) {
		form := new(RegisterForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := new(models.User)
		utils.Copy(form, u)
		if form.Password != nil {
			password, _ := goaccount.HashPassword(*form.Password)
			u.Password = &password
		}

		if form.Username == nil {
			u.Username = GenerateUsername(u.Email)
		}

		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := u.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		otp, err := models.NewOTP(ctx, u.ID, "AUTH")

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Couldn't save OTP",
			})
			return
		}

		//Sending Email
		items := map[string]string{"code": strconv.Itoa(otp.Code)}
		services.SendEmail(services.EmailConfig{
			Approach:    services.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "OTP Code",
			Template:    "otp",
			Args:        items,
		})

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.POST("/refresh", func(c *gin.Context) {
		form := new(RefreshTokenForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err := goaccount.VerifyToken(form.RefreshToken)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tb := models.TokenBlacklist{
			Token: form.RefreshToken,
		}
		ctx, _ := c.MustGet("ctx").(context.Context)
		if err := tb.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokens, err := goaccount.GenerateFullTokens(claims.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
	})

	g.POST("/otp", func(c *gin.Context) {
		form := new(OTPSendForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "User does not found",
			})
			return
		}

		ctx, _ := c.MustGet("ctx").(context.Context)

		otp, err := models.GetOTPByUserID(u.ID)
		if err != nil {
			otp, err = models.NewOTP(ctx, u.ID, "AUTH")
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   err.Error(),
					"message": "Couldn't save OTP",
				})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Threre's still a valid OTP Code try to resend it",
				"message": "Couldn't save OTP",
			})
			return
		}

		//Sending Email
		items := map[string]string{"code": strconv.Itoa(otp.Code)}
		if u.FirstName != nil {
			items["name"] = *u.FirstName
		}

		services.SendEmail(services.EmailConfig{
			Approach:    services.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "OTP Code",
			Template:    "otp",
			Args:        items,
		})

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.POST("/otp/resend", func(c *gin.Context) {
		form := new(OTPSendForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "User does not found",
			})
			return
		}

		ctx, _ := c.MustGet("ctx").(context.Context)

		otp, err := models.GetOTPByUserID(u.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Code doesn't exists try to create it first",
			})
			return
		}

		if time.Now().Before(otp.SentAt.Add(2 * time.Minute)) {
			timeRemaining := otp.SentAt.Add(2 * time.Minute).Sub(time.Now()).Round(1 * time.Second)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Retry timeout",
				"message": fmt.Sprintf("You should wait %s before sending another code", timeRemaining),
			})
			return
		} else {
			otp.UpdateSentAt(ctx)
		}

		//Sending Email
		items := map[string]string{"code": strconv.Itoa(otp.Code)}
		if u.FirstName != nil {
			items["name"] = *u.FirstName
		}

		services.SendEmail(services.EmailConfig{
			Approach:    services.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "OTP Code",
			Template:    "otp",
			Args:        items,
		})

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.POST("/otp/verify", func(c *gin.Context) {
		form := new(OTPConfirmForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}

		//Verifying OTP
		ctx, _ := c.MustGet("ctx").(context.Context)
		otp := models.OTP{
			UserID: u.ID,
			Code:   form.Code,
		}

		err = otp.Verify(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "A problem occured when trying to verify the code",
			})
			return
		}
		if !otp.IsVerified {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   nil,
				"message": "Code does not found or it is wrong",
			})
			return
		}

		//Verifying User
		u.Status = "ACTIVE"
		if err := u.Verify(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if otp.Perpose == "FORGET_PASSWORD" {
			if err := u.ExpirePassword(ctx); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		//Generating Token
		tokens, err := goaccount.GenerateFullTokens(u.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tokens)
	})

	g.POST("/password/forget", func(c *gin.Context) {

		form := new(OTPSendForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "User does not found",
			})
			return
		}

		//Creating OTP
		ctx, _ := c.MustGet("ctx").(context.Context)
		otp, err := models.NewOTP(ctx, u.ID, "FORGET_PASSWORD")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Couldn't save OTP",
			})
			return
		}

		//Sending Email
		items := map[string]string{"code": strconv.Itoa(otp.Code)}
		if u.FirstName != nil {
			items["name"] = *u.FirstName
		}

		services.SendEmail(services.EmailConfig{
			Approach:    services.EmailApproachTemplate,
			Destination: u.Email,
			Title:       "Forget Password OTP Code",
			Template:    "forget-password",
			Args:        items,
		})

		c.JSON(http.StatusOK, gin.H{})

	})

	g.PUT("/password", LoginRequired(), func(c *gin.Context) {

		ctx, _ := c.MustGet("ctx").(context.Context)
		user := c.MustGet("user").(*models.User)
		var password string

		if user.PasswordExpired || user.Password == nil {

			//Direct Password change
			form := new(DirectPasswordChangeForm)
			if err := c.ShouldBindJSON(form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			password = form.Password

		} else {

			//Normal Password change
			form := new(NormalPasswordChangeForm)
			if err := c.ShouldBindJSON(form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := goaccount.CheckPasswordHash(form.CurrentPassword, *user.Password); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
				return
			}
			password = form.Password
		}

		newPassword, err := goaccount.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = &newPassword
		if err := user.UpdatePassword(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "success"})

	})

	g.POST("/pre-register", func(c *gin.Context) {

		form := new(PreRegisterForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		emailStatus := "UNKOWN"
		usernameStatus := "UNKOWN"

		if form.Email != nil {
			u, err := models.GetUserByEmail(*form.Email)
			emailStatus = "AVAILABLE"
			if err == nil && u.Status == "ACTIVE" {
				emailStatus = "EXISTS"
			}
		}
		if form.Username != nil {
			u, err := models.GetUserByUsername(*form.Username)
			usernameStatus = "AVAILABLE"
			if err == nil && u.Status == "ACTIVE" {
				usernameStatus = "EXISTS"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"email":    emailStatus,
			"username": usernameStatus,
		})

	})
}

func GenerateUsername(email string) string {
	var username string = email
	var re *regexp.Regexp

	re = regexp.MustCompile("@.*$")
	username = re.ReplaceAllString(username, "")

	re = regexp.MustCompile("[^a-z0-9._-]")
	username = re.ReplaceAllString(username, "-")

	re = regexp.MustCompile("[._-]{2,}")
	username = re.ReplaceAllString(username, "-")

	username = strings.ToLower(username)
	username = username[0:int(math.Min(float64(len(username)), 20))]

	username = username + strconv.Itoa(int(1000+rand.Float64()*9000))

	return username
}
