package models

import (
	"context"
	"errors"
	"net/url"
	"shin/src/config"
	"shin/src/database"
	"shin/src/wallet"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Credential struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`

	SchemaID uuid.UUID `db:"schema_id" json:"schema_id"`
	Schema   *Schema   `db:"-" json:"schema"`

	CreatedID uuid.UUID `db:"created_id" json:"created_id"`
	Created   *User     `db:"-" json:"created"`

	OrganizationID uuid.UUID     `db:"organization_id" json:"organization_id"`
	Organization   *Organization `db:"-" json:"organization"`

	RecipientID *uuid.UUID `db:"recipient_id" json:"recipient_id"`
	Recipient   *Recipient `db:"-" json:"recipient"`

	ConnectionID  *string        `db:"connection_id" json:"connection_id"`
	ConnectionURL *string        `db:"connection_url" json:"connection_url"`
	Claims        types.JSONText `db:"claims" json:"claims"`

	Status CredentialStatusType `db:"status" json:"status"`

	ConnectionAt *time.Time `db:"connection_at" json:"connection_at"`
	IssuedAt     *time.Time `db:"issued_at" json:"issued_at"`
	ExpiredAt    *time.Time `db:"expired_at" json:"expired_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`

	RecipientJson    types.JSONText `db:"recipient" json:"-"`
	CreatedJson      types.JSONText `db:"created" json:"-"`
	SchemaJson       types.JSONText `db:"schema" json:"-"`
	OrganizationJson types.JSONText `db:"organization" json:"-"`
}

func (Credential) TableName() string {
	return "credentials"
}

func (Credential) FetchQuery() string {
	return "credentials/fetch"
}

func (c *Credential) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/create",
		c.Name, c.Description, c.SchemaID, c.CreatedID, c.OrganizationID, c.RecipientID, c.Claims,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}
	return database.Fetch(c, c.ID)
}

func (v *Credential) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/update",
		v.ID, v.Name, v.Description, v.SchemaID, v.Claims,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(v, v.ID)
}

func (c *Credential) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "credentials/delete", c.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (c *Credential) NewConnection(ctx context.Context, callback string) error {
	conn, err := wallet.CreateConnection(callback)
	if err != nil {
		return err
	}
	connectURL, _ := url.JoinPath(config.Config.Host, conn.ShortID)
	rows, err := database.Query(
		ctx,
		"credentials/update_connection",
		c.ID, conn.ID, connectURL,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return database.Fetch(c, c.ID)
}

func (c *Credential) Issue(ctx context.Context) error {
	if c.ConnectionID == nil {
		return errors.New("connection not valid")
	}
	if time.Since(*c.ConnectionAt) > time.Hour {
		return errors.New("connection expired")
	}

	if err := c.Organization.NewDID(ctx); err != nil {
		return err
	}
	if _, err := wallet.SendCredential(*c.ConnectionID, *c.Organization.DID, c.Claims); err != nil {
		return err
	}
	rows, err := database.Query(
		ctx,
		"credentials/update_issuing",
		c.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetCredential(id uuid.UUID) (*Credential, error) {
	c := new(Credential)

	if err := database.Fetch(c, id); err != nil {
		return nil, err
	}
	return c, nil
}

func GetCredentials(userId uuid.UUID, p database.Paginate) ([]Credential, int, error) {
	var (
		credentials = []Credential{}
		fetchList   []database.FetchList
		ids         []interface{}
	)

	if err := database.QuerySelect("credentials/get", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return credentials, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&credentials, ids...); err != nil {
		return nil, 0, err
	}
	return credentials, fetchList[0].TotalCount, nil
}
