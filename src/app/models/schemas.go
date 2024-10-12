package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Schema struct {
	ID            uuid.UUID   `db:"id" json:"id"`
	Name          string      `db:"name" json:"name"`
	Description   *string     `db:"description" json:"description"`
	CreatedID     *uuid.UUID  `db:"created_id" json:"created_id"`
	Created       *User       `db:"-" json:"created"`
	Public        bool        `db:"public" json:"public"`
	IssueDisabled bool        `db:"issue_disabled" json:"issue_disabled"`
	Deleteable    bool        `db:"deleteable" json:"deleteable"`
	Attributes    []Attribute `db:"-" json:"attributes"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`

	AttributesJson types.JSONText `db:"attributes" json:"-"`
	CreatedJson    types.JSONText `db:"created" json:"-"`
}

type Attribute struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	Name        string        `db:"name" json:"name"`
	Description *string       `db:"description" json:"description"`
	SchemaID    uuid.UUID     `db:"schema_id" json:"-"`
	Type        AttributeType `db:"type" json:"type"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
}

func (Schema) TableName() string {
	return "credential_schemas"
}

func (Schema) FetchQuery() string {
	return "credentials/fetch_schema"
}

func (s *Schema) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "credentials/delete_schema", s.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (s *Schema) Create(ctx context.Context) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(
		ctx,
		tx,
		"credentials/create_schema",
		s.Name, s.Description, s.CreatedID, s.Public,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	for rows.Next() {
		if err := rows.StructScan(s); err != nil {
			tx.Rollback()
			return err
		}
	}
	rows.Close()
	for i := range s.Attributes {
		s.Attributes[i].SchemaID = s.ID
	}

	if _, err := database.TxExecuteQuery(tx, "credentials/create_attributes", s.Attributes); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return database.Fetch(s, s.ID)
}

func GetSchema(id uuid.UUID) (*Schema, error) {
	s := new(Schema)

	if err := database.Fetch(s, id); err != nil {
		return nil, err
	}
	return s, nil
}

func GetSchemas(userId uuid.UUID, p database.Paginate) ([]Schema, int, error) {
	var (
		schemas   = []Schema{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("credentials/get_schemas", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return schemas, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&schemas, ids...); err != nil {
		return nil, 0, err
	}
	return schemas, fetchList[0].TotalCount, nil
}
