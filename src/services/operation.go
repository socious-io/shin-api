package services

import (
	"context"
	"fmt"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/utils"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

var OperationChannel = CategorizeChannel("operation")

type OperationConfig struct {
	Trigger models.OperationServiceTrigger
	Entity  interface{}
}

func SendOperation(importConfig OperationConfig) {
	Mq.sendJson(OperationChannel, importConfig)
}

func OperationWorker(message interface{}) {
	config := new(OperationConfig)
	utils.Copy(message, config)

	var (
		trigger = config.Trigger
		entity  = config.Entity
	)

	if trigger == models.OperationCredentialRevoke {
		credentialRevoke(entity)
		return
	}

}

/*
	Operation Definitions
*/

func credentialRevoke(entity interface{}) {

	ctx := context.Background()
	credential := new(models.Credential)
	utils.Copy(entity, credential)

	if err := credential.Revoke(ctx); err != nil {
		fmt.Println("Couldn't Revoke Credential, Error: ", err.Error())
		return
	}
	return

}

func CredentialBulkEmailAsync(ctx context.Context, userId uuid.UUID, schemaId string, message string, doneChan chan bool, errChan chan error) {

	p := database.Paginate{
		Limit: 100,
		Offet: 0,
		Filters: []database.Filter{
			{
				Key:   "schema_id",
				Value: schemaId,
			},
			{
				Key:   "sent",
				Value: "false",
			},
		},
	}

	for {
		credentials, total, err := models.GetCredentials(userId, p)
		ids := []uuid.UUID{}

		if err != nil {
			errChan <- err
			return
		} else if total < 1 {
			break
		}

		for _, credential := range credentials {

			fmt.Printf("Sending: %s", credential.Name)
			if credential.Recipient != nil && credential.Recipient.Email != nil {
				items := map[string]string{
					"title":      credential.Name,
					"issuer_org": credential.Organization.Name,
					"recipient":  fmt.Sprintf("%s %s", *credential.Recipient.FirstName, *credential.Recipient.LastName),
					"link":       fmt.Sprintf("%s/connect/credential/%s", config.Config.FrontHost, credential.ID.String()),
					"message":    message,
				}
				SendEmail(EmailConfig{
					Approach:    EmailApproachTemplate,
					Destination: *credential.Recipient.Email,
					Title:       "Shin: Your Verifiable Credential is Here",
					Template:    "credentials-recipients",
					Args:        items,
				})
			}

			ids = append(ids, credential.ID)

		}

		p.Offet = p.Offet + p.Limit
		models.CredentialsBulkSend(ctx, ids, userId)
		ids = []uuid.UUID{}

	}

	doneChan <- true
	return

}
