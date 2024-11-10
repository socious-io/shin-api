package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type CSVImport struct {
	ID        uuid.UUID        `db:"id" json:"id"`
	DocType   CSVImportDocType `db:"doc_type" json:"doc_type"`
	UserID    uuid.UUID        `db:"user_id" json:"user_id"`
	Data      types.JSONText   `db:"data" json:"data"`
	Status    CSVImportStatus  `db:"status" json:"status"`
	Reason    *string          `db:"reason" json:"reason"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
}

func (CSVImport) TableName() string {
	return "csv_imports"
}

func (CSVImport) FetchQuery() string {
	return "csv_imports/fetch"
}

func (ci *CSVImport) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(ci)
}

func (ci *CSVImport) Create(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"csv_imports/create",
		ci.UserID, ci.DocType, ci.Data, ci.Status, ci.Reason,
	)

	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := ci.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func GetCSVImport(id uuid.UUID) (*CSVImport, error) {
	ci := new(CSVImport)
	if err := database.Fetch(ci, id); err != nil {
		return nil, err
	}
	return ci, nil
}
