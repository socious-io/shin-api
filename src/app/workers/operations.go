package workers

import (
	"context"
	"fmt"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/utils"

	"github.com/socious-io/gomail"
	"github.com/socious-io/gomq"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

var OperationChannel = "operation" //TODO: Categorize

type OperationParams struct {
	Trigger models.OperationServiceTrigger
	Entity  interface{}
}

func SendOperation(params OperationParams) {
	gomq.Mq.SendJson(OperationChannel, params)
}

func OperationWorker(params OperationParams) error {
	var (
		trigger = params.Trigger
		entity  = params.Entity
	)

	switch trigger {
	case models.OperationCredentialRevoke:
		return credentialRevoke(entity)

	}

	return nil
}

/*
	Operation Definitions
*/

func credentialRevoke(entity interface{}) error {

	ctx := context.Background()
	credential := new(models.Credential)
	utils.Copy(entity, credential)

	if err := credential.Revoke(ctx); err != nil {
		fmt.Println("Couldn't Revoke Credential, Error: ", err.Error())
		return err
	}
	return nil

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
			if credential.Recipient != nil && credential.Recipient.Email != nil {
				gomail.SendEmail(gomail.EmailConfig{
					Approach:    gomail.EmailApproachTemplate,
					Destination: *credential.Recipient.Email,
					Title:       "Shin: Your Verifiable Credential is Here",
					TemplateId:  "credentials-recipients",
					Args: map[string]string{
						"title":      credential.Name,
						"issuer_org": credential.Organization.Name,
						"recipient":  fmt.Sprintf("%s %s", *credential.Recipient.FirstName, *credential.Recipient.LastName),
						"link":       fmt.Sprintf("%s/connect/credential/%s", config.Config.FrontHost, credential.ID.String()),
						"message":    message,
					},
				})
			}

			ids = append(ids, credential.ID)
		}

		p.Offet = p.Offet + p.Limit
		models.CredentialsBulkSend(ctx, ids, userId)
	}

	doneChan <- true
}
