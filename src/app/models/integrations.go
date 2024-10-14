package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IntegrationKey struct {
	ID        uuid.UUID                 `db:"id" json:"id"`
	UserID    uuid.UUID                 `db:"user_id" json:"-"`
	Name      string                    `db:"name" json:"name"`
	Key       string                    `db:"key" json:"key"`
	Secret    string                    `db:"secret" json:"secret"`
	BaseUrl   string                    `db:"base_url" json:"base_url"`
	Status    KybVerificationStatusType `db:"status" json:"status"`
	CreatedAt time.Time                 `db:"created_at" json:"created_at"`
}

func (IntegrationKey) TableName() string {
	return "integration_keys"
}

func (IntegrationKey) FetchQuery() string {
	return "integrations/fetch_key"
}

func (ik *IntegrationKey) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(ik)
}

func (ik *IntegrationKey) Create(ctx context.Context) (*IntegrationKey, error) {

	rows, err := database.Query(
		ctx,
		"integrations/create_key",
		ik.Name, ik.UserID, ik.BaseUrl, ik.Key, ik.Secret,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := ik.Scan(rows); err != nil {
			return nil, err
		}
	}
	return ik, nil
}

func (ik *IntegrationKey) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"integrations/update_key",
		ik.ID, ik.Name,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(ik); err != nil {
			return err
		}
	}
	return nil
}

func (ik *IntegrationKey) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "integrations/delete_key", ik.ID, ik.UserID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetKeyById(id uuid.UUID) (*IntegrationKey, error) {
	k := new(IntegrationKey)
	if err := database.Get(k, "integrations/fetch_key_by_id", id); err != nil {
		return nil, err
	}
	return k, nil
}

func GetAllKeysByUserId(userId uuid.UUID, p database.Paginate) ([]IntegrationKey, int, error) {

	var (
		integrationKeys = []IntegrationKey{}
		fetchList       []database.FetchList
		ids             []interface{}
	)

	if err := database.QuerySelect("integrations/fetch_keys_by_userid", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return integrationKeys, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&integrationKeys, ids...); err != nil {
		return nil, 0, err
	}
	return integrationKeys, fetchList[0].TotalCount, nil
}
