package views

import (
	"shin/src/app/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

// Authentication
type RegisterForm struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"username"`
	Email     string  `json:"email" validate:"required,email"`
	Password  *string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type OTPSendForm struct {
	Email string `json:"email" validate:"required,email"`
}
type OTPConfirmForm struct {
	Email string `json:"email" validate:"required,email"`
	Code  int    `json:"code" validate:"required"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type PreRegisterForm struct {
	Email    *string `json:"email" validate:"email"`
	Username *string `json:"username"`
}

type NormalPasswordChangeForm struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type DirectPasswordChangeForm struct {
	Password string `json:"password" validate:"required"`
}

type AuthForm struct {
	RedirectURL string `json:"redirect_url" validate:"required"`
}

type SessionForm struct {
	Code string `json:"code" validate:"required"`
}

type SyncForm struct {
	Organizations []models.Organization `json:"organizations"`
	User          models.User           `json:"user" validate:"required"`
}

// Others
type OrganizationForm struct {
	Name        string          `json:"name" validate:"required,min=3,max=32"`
	Description string          `json:"description" validate:"required,min=3"`
	Logo        *types.JSONText `json:"logo"`
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

type CredentialBySchemaEmailForm struct {
	SchemaID uuid.UUID `json:"schema_id" validate:"required"`
	Message  string    `json:"message" validate:"required"`
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
	Username  *string         `json:"username" validate:"required,min=3,max=32"`
	JobTitle  *string         `json:"job_title"`
	Bio       *string         `json:"bio"`
	FirstName string          `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string          `json:"last_name" validate:"required,min=3,max=32"`
	Phone     *string         `json:"phone"`
	Avatar    *types.JSONText `json:"avatar"`
}

type KYBVerificationForm struct {
	Documents []string `json:"documents"`
}

type ApikeyForm struct {
	Name string `json:"name"`
}
