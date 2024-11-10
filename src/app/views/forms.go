package views

import (
	"shin/src/app/models"

	"github.com/google/uuid"
)

type OrganizationForm struct {
	Name        string     `json:"name" validate:"required,min=3,max=32"`
	Description string     `json:"description" validate:"required,min=3"`
	LogoID      *uuid.UUID `json:"logo_id"`
}

type SchemaForm struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Public      bool    `json:"public"`
	Attributes  []struct {
		Name        string               `json:"name"`
		Description *string              `json:"description"`
		Type        models.AttributeType `json:"type"`
	} `json:"attributes"`
}

type VerificationForm struct {
	Name        string                  `json:"name" validate:"required,min=3,max=32"`
	Description *string                 `json:"description" validate:"required,min=3"`
	SchemaID    uuid.UUID               `json:"schema_id" validate:"required"`
	Type        models.VerificationType `json:"type" validate:"required"`
	Attributes  []struct {
		AttributeID uuid.UUID                       `json:"attribute_id"`
		Operator    models.VerificationOperatorType `json:"operator"`
		Value       string                          `json:"value"`
	} `json:"attributes"`
}

type VerificationIndividualForm struct {
	CustomerID     string    `json:"customer_id" validate:"required"`
	VerificationID uuid.UUID `json:"verification_id" validate:"required"`
}

type Claim struct {
	Name  string      `json:"name" validate:"required,min=3,max=32"`
	Value interface{} `json:"value" validate:"required"`
}

type CredentialForm struct {
	Name        string    `json:"name" validate:"required,min=3,max=32"`
	Description *string   `json:"description" validate:"required,min=3"`
	SchemaID    uuid.UUID `json:"schema_id" validate:"required"`
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
	Claims      []Claim   `json:"claims" validate:"required"`
}

type CredentialBulkOperationForm struct {
	Credentials []uuid.UUID `json:"credentials" validate:"required"`
}

type CredentialBulkEmailForm struct {
	CredentialBulkOperationForm
	Message string `json:"message" validate:"required"`
}

type RecipientForm struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=128"`
	LastName  string `json:"last_name" validate:"required,min=3,max=128"`
	Email     string `json:"email" validate:"required,email"`
}

type CredentialRecipientForm struct {
	Credential CredentialForm `json:"credential"`
	Recipient  RecipientForm  `json:"recipient"`
}

type ProfileUpdateForm struct {
	Username  *string    `json:"username" validate:"required,min=3,max=32"`
	JobTitle  *string    `json:"job_title"`
	Bio       *string    `json:"bio"`
	FirstName string     `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string     `json:"last_name" validate:"required,min=3,max=32"`
	Phone     *string    `json:"phone"`
	AvatarID  *uuid.UUID `json:"avatar_id"`
}

type KYBVerificationForm struct {
	Documents []string `json:"documents"`
}

type ApikeyForm struct {
	Name string `json:"name"`
}
