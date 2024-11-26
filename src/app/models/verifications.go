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
	ID          uuid.UUID               `db:"id" json:"id"`
	Name        string                  `db:"name" json:"name"`
	Description *string                 `db:"description" json:"description"`
	SchemaID    uuid.UUID               `db:"schema_id" json:"schema_id"`
	Schema      *Schema                 `db:"-" json:"schema"`
	UserID      uuid.UUID               `db:"user_id" json:"user_id"`
	User        *User                   `db:"-" json:"user"`
	Attributes  []VerificationAttribute `db:"-" json:"attributes"`
	Type        VerificationType        `db:"type" json:"type"`
	Single      *VerificationIndividual `db:"-" json:"single"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	AttributesJson types.JSONText `db:"attributes" json:"-"`
	UserJson       types.JSONText `db:"user" json:"-"`
	SchemaJson     types.JSONText `db:"schema" json:"-"`
	SingleJson     types.JSONText `db:"single" json:"-"`
}

type VerificationIndividual struct {
	ID uuid.UUID `db:"id" json:"id"`

	UserID uuid.UUID `db:"user_id" json:"user_id"`
	User   *User     `db:"-" json:"user"`

	RecipientID uuid.UUID  `db:"recipient_id" json:"recipient_id"`
	Recipient   *Recipient `db:"-" json:"recipient"`

	VerificationID uuid.UUID     `db:"verification_id" json:"verification_id"`
	Verification   *Verification `db:"-" json:"verification"`

	PresentID       *string                `db:"present_id" json:"present_id"`
	ConnectionID    *string                `db:"connection_id" json:"connection_id"`
	ConnectionURL   *string                `db:"connection_url" json:"connection_url"`
	Body            types.JSONText         `db:"body" json:"body"`
	Status          VerificationStatusType `db:"status" json:"status"`
	ValidationError *string                `db:"validation_error" json:"validation_error"`

	ConnectionAt     *time.Time     `db:"connection_at" json:"connection_at"`
	VerifiedAt       *time.Time     `db:"verified_at" json:"verified_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UserJson         types.JSONText `db:"user" json:"-"`
	VerificationJson types.JSONText `db:"verification" json:"-"`
	RecipientJson    types.JSONText `db:"recipient" json:"-"`
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
	return "credential_verificationss"
}

func (Verification) FetchQuery() string {
	return "verifications/fetch"
}

func (VerificationIndividual) TableName() string {
	return "credential_individuals"
}

func (VerificationIndividual) FetchQuery() string {
	return "verifications/fetch_individual"
}

func (v *Verification) Create(ctx context.Context) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(
		ctx,
		tx,
		"verifications/create",
		v.Name, v.Description, v.UserID, v.SchemaID, v.Type,
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

	for i := range v.Attributes {
		v.Attributes[i].VerificationID = v.ID
		v.Attributes[i].SchemaID = v.SchemaID
	}
	if len(v.Attributes) > 0 {
		if _, err := database.TxExecuteQuery(tx, "verifications/create_attributes", v.Attributes); err != nil {
			tx.Rollback()
			return err
		}
	}
	if v.Type == VerificationSingle {
		noneCustomerID := "Any"

		r := &Recipient{
			UserID:     v.UserID,
			CustomerID: &noneCustomerID,
		}
		if err := r.Create(ctx); err != nil {
			tx.Rollback()
			return err
		}
		rows, err = database.TxQuery(
			ctx,
			tx,
			"verifications/create_individual",
			v.UserID, r.ID, v.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows.Close()
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
		"verifications/update",
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

	rows, err = database.TxQuery(ctx, tx, "verifications/delete_attributes", v.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	for i := range v.Attributes {
		v.Attributes[i].VerificationID = v.ID
		v.Attributes[i].SchemaID = v.SchemaID
	}
	if len(v.Attributes) > 0 {
		if _, err := database.TxExecuteQuery(tx, "verifications/create_attributes", v.Attributes); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return database.Fetch(v, v.ID)
}

func (v *Verification) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "verifications/delete", v.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (v *VerificationIndividual) Create(ctx context.Context, customerID string) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	r := &Recipient{
		UserID:     v.UserID,
		CustomerID: &customerID,
	}
	if err := r.Create(ctx); err != nil {
		tx.Rollback()
		return err
	}
	v.RecipientID = r.ID
	rows, err := database.TxQuery(
		ctx,
		tx,
		"verifications/create_individual",
		v.UserID, r.ID, v.VerificationID,
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

	if err := tx.Commit(); err != nil {
		return err
	}
	return database.Fetch(v, v.ID)
}

func (v *VerificationIndividual) NewConnection(ctx context.Context, callback string) error {
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
		"verifications/update_connection",
		v.ID, conn.ID, connectURL,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (v *VerificationIndividual) ProofRequest(ctx context.Context) error {
	if v.ConnectionID == nil {
		return errors.New("connection not valid")
	}
	if time.Since(*v.ConnectionAt) > time.Hour {
		return errors.New("connection expired")
	}

	challenge, _ := json.Marshal(wallet.H{
		"type": v.Verification.Schema.Name,
	})

	presentID, err := wallet.ProofRequest(*v.ConnectionID, string(challenge))
	if err != nil {
		return err
	}
	rows, err := database.Query(
		ctx,
		"verifications/update_present_id",
		v.ID, presentID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (v *VerificationIndividual) ProofVerify(ctx context.Context) error {
	if v.PresentID == nil {
		return errors.New("need request proof present first")
	}

	vc, err := wallet.ProofVerify(*v.PresentID)
	if err != nil {
		return err
	}
	vcData, _ := json.Marshal(vc)
	if len(v.Verification.Attributes) > 0 {
		if err := validateVC(*v.Verification.Schema, vc, v.Verification.Attributes); err != nil {
			rows, err := database.Query(
				ctx,
				"verifications/update_present_failed",
				v.ID, vcData, err.Error(),
			)
			if err != nil {
				return err
			}
			rows.Close()
			return nil
		}
	}
	rows, err := database.Query(
		ctx,
		"verifications/update_present_verify",
		v.ID, vcData,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func GetVerification(id uuid.UUID) (*Verification, error) {
	v := new(Verification)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	schema, err := GetSchema(v.SchemaID)
	if err != nil {
		return nil, err
	}
	v.Schema = schema
	return v, nil
}

func GetVerifications(userId uuid.UUID, p database.Paginate) ([]Verification, int, error) {
	var (
		verifications = []Verification{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if len(p.Filters) > 0 && p.Filters[0].Key == "type" {
		if err := database.QuerySelect("verifications/get_by_type", &fetchList, userId, p.Limit, p.Offet, p.Filters[0].Value); err != nil {
			return nil, 0, err
		}
	} else {
		if err := database.QuerySelect("verifications/get", &fetchList, userId, p.Limit, p.Offet); err != nil {
			return nil, 0, err
		}
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

func GetVerificationsIndividual(id uuid.UUID) (*VerificationIndividual, error) {
	v := new(VerificationIndividual)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	verification, err := GetVerification(v.VerificationID)
	if err != nil {
		return nil, err
	}
	v.Verification = verification
	return v, nil
}

func GetVerificationsIndividuals(userId, verificationId uuid.UUID, p database.Paginate) ([]VerificationIndividual, int, error) {
	var (
		verifications = []VerificationIndividual{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("verifications/get_individuals", &fetchList, userId, verificationId, p.Limit, p.Offet); err != nil {
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
	if err := database.Fetch(&schema, schema.ID); err != nil {
		return err
	}

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
			if val <= attrVal {
				return validationErr
			}
		case OperatorSmaller:
			val, attrVal, err := convertValsToNumber(value, attr.Value)
			if err != nil {
				return err
			}
			if val >= attrVal {
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
		customDateLayout = "2006-01-02"
		val              int
		isTime           bool = false
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
			if t, err := time.Parse(customDateLayout, v); err == nil {
				val = int(t.Unix())
				isTime = true
			}
		}
	case int:
		val = v
	}
	if isTime {
		if t, err := time.Parse(time.RFC3339, attrVal); err == nil {
			return val, int(t.Unix()), nil
		}
	}
	attrIntVal, err := strconv.Atoi(attrVal)
	if err != nil {
		return 0, 0, fmt.Errorf("could not operate bigger/smaller on not number/date values")
	}
	return val, attrIntVal, nil
}
