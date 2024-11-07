package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"shin/src/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ImportCredentials(recordId uuid.UUID, meta map[string]any) error {

	//check if user has orgs

	//TODO: should we import:
	// Name        string    `json:"name" validate:"required,min=3,max=32"` // Credential name
	// Description *string   `json:"description" validate:"required,min=3"` // Credential description
	schema_id, user_id := uuid.MustParse(meta["schema_id"].(string)), uuid.MustParse(meta["user_id"].(string))

	s, _ := models.GetSchema(schema_id)
	u, _ := models.GetUser(user_id)
	ctx := context.Background()

	ci, _ := models.GetCSVImport(recordId)
	var import_error error = nil

	var fields []map[string]string
	json.Unmarshal(ci.Data, &fields)

	for _, record := range fields {

		//Extract recipient info and create it
		FirstName, LastName, Email := record["first_name"], record["last_name"], record["email"]
		r := models.Recipient{
			FirstName: &FirstName,
			LastName:  &LastName,
			Email:     &Email,
			UserID:    u.ID,
		}
		r.UserID = u.ID
		if err := r.Create(ctx.(context.Context)); err != nil {
			import_error = err
			break
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
			import_error = fmt.Errorf("fetching org error :%v", err)
			break
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
			import_error = err
			break
		}
	}

	if import_error != nil {
		err_str := import_error.Error()
		ci.Status = models.CSVStatusValidationFailed
		ci.Reason = &err_str
	} else {
		ci.Status = models.CSVStatusDone
	}

	ci.Update(ctx.(context.Context))
	return import_error

}
