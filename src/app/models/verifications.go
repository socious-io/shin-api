package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"shin/src/config"
	"shin/src/database"
	"shin/src/wallet"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Verification struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	SchemaID    uuid.UUID `db:"schema_id" json:"schema_id"`
	Schema      *Schema   `db:"-" json:"schema"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	User        *User     `db:"-" json:"user"`

	PresentID     *string                 `db:"present_id" json:"present_id"`
	ConnectionID  *string                 `db:"connection_id" json:"connection_id"`
	ConnectionURL *string                 `db:"connection_url" json:"connection_url"`
	Body          types.JSONText          `db:"body" json:"body"`
	Attributes    []VerificationAttribute `db:"-" json:"attributes"`
	Status        VerificationStatusType  `db:"status" json:"status"`

	ConnectionAt *time.Time `db:"connection_at" json:"connection_at"`
	VerifiedAt   *time.Time `db:"verified_at" json:"verified_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`

	AttributesJson types.JSONText `db:"attributes" json:"-"`
	UserJson       types.JSONText `db:"user" json:"-"`
	SchemaJson     types.JSONText `db:"schema" json:"-"`
}

type VerificationAttribute struct {
	ID             uuid.UUID                `db:"id" json:"id"`
	AttributeID    uuid.UUID                `db:"attribute_id" json:"attribute_id"`
	SchemaID       uuid.UUID                `db:"schema_id" json:"schema_id"`
	VerificationID uuid.UUID                `db:"verification_id" json:"verification_id"`
	Value          string                   `db:"value" json:"value"`
	Operator       VerificationOperatorType `db:"operator" json:"operator"`
	CreatedAt      time.Time                `db:"created_at" json:"created_at"`
}

func (Verification) TableName() string {
	return "credential_verifications"
}

func (Verification) FetchQuery() string {
	return "credentials/fetch_verification"
}

func (v *Verification) Create(ctx context.Context) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(
		ctx,
		tx,
		"credentials/create_verification",
		v.Name, v.Description, v.UserID, v.SchemaID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := range v.Attributes {
		v.Attributes[i].VerificationID = v.ID
		v.Attributes[i].SchemaID = v.SchemaID
	}
	if len(v.Attributes) > 0 {
		if _, err := database.TxExecuteQuery(tx, "credentials/create_verification_attributes", v.Attributes); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return database.Fetch(v, v.ID)
}

func (v *Verification) Update(ctx context.Context) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(
		ctx,
		tx,
		"credentials/update_verification",
		v.ID, v.Name, v.Description, v.UserID, v.SchemaID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	rows, err = database.TxQuery(ctx, tx, "credentials/delete_verification_attributes", v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	for i := range v.Attributes {
		v.Attributes[i].VerificationID = v.ID
		v.Attributes[i].SchemaID = v.SchemaID
	}

	if _, err := database.TxExecuteQuery(tx, "credentials/create_verification_attributes", v.Attributes); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return database.Fetch(v, v.ID)
}

func (v *Verification) NewConnection(ctx context.Context, callback string) error {
	if v.Status == StatusVerifRequested {
		return nil
	}
	conn, err := wallet.CreateConnection(callback)
	if err != nil {
		return err
	}
	connectURL, _ := url.JoinPath(config.Config.Host, conn.ShortID)
	rows, err := database.Query(
		ctx,
		"credentials/update_connection_verification",
		v.ID, conn.ID, connectURL,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (v *Verification) ProofRequest(ctx context.Context) error {
	if v.ConnectionID == nil {
		return errors.New("connection not valid")
	}
	if time.Since(*v.ConnectionAt) > time.Hour {
		return errors.New("connection expired")
	}

	challenge, _ := json.Marshal(wallet.H{
		"type": v.Schema.Name,
	})

	presentID, err := wallet.ProofRequest(*v.ConnectionID, string(challenge))
	if err != nil {
		return err
	}
	rows, err := database.Query(
		ctx,
		"credentials/update_present_id_verification",
		v.ID, presentID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (v *Verification) ProofVerify(ctx context.Context) error {
	if v.PresentID == nil {
		return errors.New("need request proof present first")
	}

	vc, err := wallet.ProofVerify(*v.PresentID)
	if err != nil {
		return err
	}
	vcData, _ := json.Marshal(vc)
	query := "credentials/update_present_verify_verification"
	if err := validateVC(*v.Schema, vc, v.Attributes); err != nil {
		query = "credentials/update_present_failed_verification"
	}
	rows, err := database.Query(
		ctx,
		query,
		v.ID, vcData,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (v *Verification) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "credentials/delete_verification", v.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetVerification(id uuid.UUID) (*Verification, error) {
	v := new(Verification)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerificationByConnection(connectionID uuid.UUID) (*Verification, error) {
	v := new(Verification)

	if err := database.Get(v, "credentials/verification_by_connection", connectionID); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerifications(userId uuid.UUID, p database.Paginate) ([]Verification, int, error) {
	var (
		verifications = []Verification{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("credentials/get_verifications", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return verifications, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&verifications, ids...); err != nil {
		return nil, 0, err
	}
	return verifications, fetchList[0].TotalCount, nil
}

func validateVC(schema Schema, vc wallet.H, attrs []VerificationAttribute) error {
	for _, attr := range attrs {
		attrName := ""
		for _, a := range schema.Attributes {
			if a.ID == attr.AttributeID {
				attrName = a.Name
				break
			}
		}
		value, ok := vc[attrName]
		if !ok {
			return fmt.Errorf("could not find expecting attribute %s", attrName)
		}

		validationErr := fmt.Errorf("validation error on %s", attrName)

		switch attr.Operator {
		case OperatorEqual:
			if fmt.Sprintf("%v", value) != attr.Value {
				return validationErr
			}
		case OperatorBigger:
			val, attrVal, err := convertValsToNumber(value, attr.Value)
			if err != nil {
				return err
			}
			if val < attrVal {
				return validationErr
			}
		case OperatorSmaller:
			val, attrVal, err := convertValsToNumber(value, attr.Value)
			if err != nil {
				return err
			}
			if val > attrVal {
				return validationErr
			}
		case OperatorNot:
			if fmt.Sprintf("%s", value) == attr.Value {
				return validationErr
			}
		}
	}
	return nil
}

func convertValsToNumber(value interface{}, attrVal string) (int, int, error) {
	var (
		val    int
		isTime bool = false
	)
	switch v := value.(type) {
	case string:
		if intVal, err := strconv.Atoi(v); err == nil {
			val = intVal
		} else {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				val = int(t.Unix())
				isTime = true
			}
		}
	case int:
		val = v
	}
	if isTime {
		if t, err := time.Parse(time.RFC3339, attrVal); err == nil {
			return int(t.Unix()), val, nil
		}
	}
	attrIntVal, err := strconv.Atoi(attrVal)
	if err != nil {
		return 0, 0, fmt.Errorf("could not operate bigger/smaller on not number/date values")
	}
	return val, attrIntVal, nil
}
