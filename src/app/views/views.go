package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	authGroup(r)
	orgGroup(r)
	userGroup(r)
	schemasGroup(r)
	rootGroup(r)
	verificationsGroup(r)
	credentialsGroup(r)
	recipientsGroup(r)
	mediaGroup(r)
	kybVerificationGroup(r)
	integrationGroup(r)
	syncGroup(r)
}
