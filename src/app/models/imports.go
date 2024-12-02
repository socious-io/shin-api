package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Import struct {
	ID         uuid.UUID    `db:"id" json:"id"`
	Target     ImportTarget `db:"target" json:"target"`
	UserID     uuid.UUID    `db:"user_id" json:"user_id"`
	Entities   []uuid.UUID  `db:"-" json:"entities"`
	TotalCount int          `db:"total_count" json:"total_count"`
	Count      int          `db:"count" json:"count"`
	Status     ImportStatus `db:"status" json:"status"`
	CreatedAt  time.Time    `db:"created_at" json:"created_at"`

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
		i.UserID, i.Target, i.TotalCount,
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

func (i *Import) Append(ctx context.Context, entity uuid.UUID) error {
	rows, err := database.Query(
		ctx,
		"imports/append",
		i.ID, entity,
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

	err := pq.Array(&i.Entities).Scan(i.EntitiesArray)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func GetActiveImportByUserId(userID uuid.UUID) (*Import, error) {
	i := new(Import)
	if err := database.Get(i, "imports/fetch_active_by_user", userID); err != nil {
		return nil, err
	}

	err := pq.Array(&i.Entities).Scan(i.EntitiesArray)
	if err != nil {
		return nil, err
	}

	return i, nil
}
