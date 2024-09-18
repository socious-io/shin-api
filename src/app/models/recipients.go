package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
)

type Recipient struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Recipient) TableName() string {
	return "recipients"
}

func (Recipient) FetchQuery() string {
	return "credentials/fetch_recipient"
}

func (r *Recipient) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/create_recipient",
		r.FirstName, r.LastName, r.Email, r.UserID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Recipient) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/update_recipient",
		r.ID, r.FirstName, r.LastName, r.Email,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(r); err != nil {
			return err
		}
	}
	return nil
}

func (r *Recipient) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "credentials/delete_recipient", r.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetRecipient(id uuid.UUID) (*Recipient, error) {
	r := new(Recipient)

	if err := database.Fetch(r, id); err != nil {
		return nil, err
	}
	return r, nil
}

func SearchRecipients(query string, userID uuid.UUID, p database.Paginate) ([]Recipient, int, error) {
	var (
		recipients = []Recipient{}
		fetchList  []database.FetchList
		ids        []interface{}
	)

	if err := database.QuerySelect(
		"credentials/search_recipients",
		&fetchList, query, userID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return recipients, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&recipients, ids...); err != nil {
		return nil, 0, err
	}
	return recipients, fetchList[0].TotalCount, nil
}
