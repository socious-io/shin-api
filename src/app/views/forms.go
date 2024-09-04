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
	Name        string    `json:"name" validate:"required,min=3,max=32"`
	Description *string   `json:"description" validate:"required,min=3"`
	SchemaID    uuid.UUID `json:"schema_id" validate:"required"`
}

type CredentialForm struct {
	Name        string    `json:"name" validate:"required,min=3,max=32"`
	Description *string   `json:"description" validate:"required,min=3"`
	SchemaID    uuid.UUID `json:"schema_id" validate:"required"`
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
	Claims      []struct {
		Name  string      `json:"name" validate:"required,min=3,max=32"`
		Value interface{} `json:"value" validate:"required"`
	} `json:"claims" validate:"required"`
}

type RecipientForm struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=128"`
	LastName  string `json:"last_name" validate:"required,min=3,max=128"`
	Email     string `json:"email" validate:"required,email"`
}

type ProfileUpdateForm struct {
	Username  string     `json:"username" validate:"required,min=3,max=32"`
	JobTitle  *string    `json:"job_title"`
	Bio       *string    `json:"bio"`
	FirstName string     `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string     `json:"last_name" validate:"required,min=3,max=32"`
	Phone     *string    `json:"phone"`
	AvatarID  *uuid.UUID `json:"avatar_id"`
}
