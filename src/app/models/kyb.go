package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type KYBDocuments struct {
	Url      string `db:"url" json:"url"`
	Filename string `db:"filename" json:"filename"`
}

type KYBVerification struct {
	ID        uuid.UUID                 `db:"id" json:"id"`
	UserID    uuid.UUID                 `db:"user_id" json:"user_id"`
	OrgID     uuid.UUID                 `db:"organization_id" json:"organization_id"`
	Status    KybVerificationStatusType `db:"status" json:"status"`
	Documents []KYBDocuments            `db:"-" json:"documents"`
	CreatedAt time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt time.Time                 `db:"updated_at" json:"updated_at"`

	//Json temp fields
	DocumentsJson types.JSONText `db:"documents" json:"-"`
}

func (KYBVerification) TableName() string {
	return "kyb_verifications"
}

func (KYBVerification) FetchQuery() string {
	return "kyb/fetch"
}

func (k *KYBVerification) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(k)
}

func (k *KYBVerification) Create(ctx context.Context, documents []string) (*KYBVerification, error) {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return nil, err
	}
	rows, err := database.TxQuery(ctx, tx, "kyb/create",
		k.UserID, k.OrgID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		if err := rows.StructScan(k); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rows.Close()

	for _, document := range documents {
		rows, err = database.TxQuery(ctx, tx, "kyb/create_document",
			k.ID, document,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		rows.Close()
	}
	tx.Commit()

	return GetById(k.ID)
}

func (k *KYBVerification) ChangeStatus(ctx context.Context, status KybVerificationStatusType) error {
	rows, err := database.Query(ctx, "kyb/change_status", k.ID, status)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := k.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func GetById(id uuid.UUID) (*KYBVerification, error) {
	k := new(KYBVerification)
	if err := database.Get(k, "kyb/fetch_by_id", id); err != nil {
		return nil, err
	}
	return k, nil
}

func GetByIdAndUserId(id, userID uuid.UUID) (*KYBVerification, error) {
	k := new(KYBVerification)
	if err := database.Get(k, "kyb/fetch_by_id_and_user", id, userID); err != nil {
		return nil, err
	}
	return k, nil
}

func GetAllByUserId(userId uuid.UUID, p database.Paginate) ([]KYBVerification, int, error) {

	var (
		kybVerifications = []KYBVerification{}
		fetchList        []database.FetchList
		ids              []interface{}
	)

	if err := database.QuerySelect("kyb/fetch_ids_by_userid", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return kybVerifications, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&kybVerifications, ids...); err != nil {
		return nil, 0, err
	}
	return kybVerifications, fetchList[0].TotalCount, nil
}
