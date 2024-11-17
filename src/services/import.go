package services

import (
	"context"
	"encoding/json"
	"fmt"
	"shin/src/app/models"
	"shin/src/config"
	"shin/src/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ImportChannel = CategorizeChannel("import")
var ImportWorkersCount = 4

type ImportConfig struct {
	Import models.Import
	Record map[string]any
	Meta   map[string]any
}

func SendImport(importConfig ImportConfig) {
	Mq.sendJson(ImportChannel, importConfig)
}

func ImportWorker(message interface{}) {
	importConfig := new(ImportConfig)
	utils.Copy(message, importConfig)

	var (
		importObj = importConfig.Import
		record    = importConfig.Record
		meta      = importConfig.Meta
	)

	if importObj.Target == models.ImportTargetCredentials {
		err := importCredentials(record, meta, importObj)
		if err != nil {
			fmt.Println("Couldn't Import Docs to Database, Error: ", err.Error())
		}
	}

}

/*
Initiate Import
*/

func InitiateImport(results []map[string]any, meta map[string]any, i models.Import) {

	for _, result := range results {
		SendImport(ImportConfig{
			Import: i,
			Record: result,
			Meta:   meta,
		})
	}
}

/*
Importing Functions ( We can have multiple import functions here )
*/
func importCredentials(record map[string]any, meta map[string]any, i models.Import) error {

	//TODO: should we import:
	// Name        string    `json:"name" validate:"required,min=3,max=32"` // Credential name
	// Description *string   `json:"description" validate:"required,min=3"` // Credential description
	schema_id, user_id := uuid.MustParse(meta["schema_id"].(string)), uuid.MustParse(meta["user_id"].(string))

	s, _ := models.GetSchema(schema_id)
	u, _ := models.GetUser(user_id)
	ctx := context.Background()

	//Extract recipient info and create it
	var FirstName, LastName, Email string = record["first_name"].(string), record["last_name"].(string), record["email"].(string)
	r := models.Recipient{
		FirstName: &FirstName,
		LastName:  &LastName,
		Email:     &Email,
		UserID:    u.ID,
	}
	r.UserID = u.ID
	if err := r.Create(ctx.(context.Context)); err != nil {
		return err
	}
	delete(record, "first_name")
	delete(record, "last_name")
	delete(record, "email")

	//Creating Credential
	cv := models.Credential{
		CreatedID:   u.ID,
		RecipientID: &r.ID,
		SchemaID:    s.ID,
	}

	//Organization
	orgs, err := models.GetOrgsByMember(cv.CreatedID)
	if err != nil || len(orgs) < 1 {
		return fmt.Errorf("fetching org error :%v", err)
	}
	cv.OrganizationID = orgs[0].ID

	//Claims
	claims := gin.H{}
	for key, claim := range record {
		claims[key] = claim
	}
	claims["type"] = s.Name
	claims["issued_date"] = time.Now().Format(time.RFC3339)
	claims["company_name"] = orgs[0].Name
	cv.Claims, _ = json.Marshal(&claims)

	if err := cv.Create(ctx.(context.Context)); err != nil {
		return err
	}

	err = i.Append(ctx, cv.ID)
	if err != nil {
		fmt.Println(err)
	}

	if i.Status == models.ImportStatusCompleted {
		SendEmail(EmailConfig{
			Approach:    EmailApproachTemplate,
			Destination: u.Email,
			Title:       "Shin: Your CSV Import has been completed", //TODO: exact title
			Template:    "credentials-import-completed",
			Args: map[string]string{
				"link": fmt.Sprintf("%s/credentials/create?schema=%s", config.Config.FrontHost, meta["schema_id"]),
			},
		})
	}

	return nil
}
