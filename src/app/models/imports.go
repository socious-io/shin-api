package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Import struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	Target    ImportTarget `db:"target" json:"target"`
	UserID    uuid.UUID    `db:"user_id" json:"user_id"`
	Entities  []uuid.UUID  `db:"-" json:"entities"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`

	EntitiesArray []uint8 `db:"entities" json:"-"`
}

func (Import) TableName() string {
	return "imports"
}

func (Import) FetchQuery() string {
	return "imports/fetch"
}

func (i *Import) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(i)
}

func (i *Import) Create(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"imports/create",
		i.UserID, i.Target,
	)

	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := i.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func (i *Import) Update(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"imports/update",
		i.ID, pq.Array(i.Entities),
	)

	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := i.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func GetImport(id uuid.UUID) (*Import, error) {
	i := new(Import)
	if err := database.Fetch(i, id); err != nil {
		return nil, err
	}
	return i, nil
}
