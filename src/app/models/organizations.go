package models

import (
	"context"
	"log"
	"shin/src/wallet"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Organization struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	DID         *string    `db:"did" json:"did"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	LogoID      *uuid.UUID `db:"logo_id" json:"logo_id"`
	Logo        struct {
		Url      *string `db:"url" json:"url"`
		Filename *string `db:"filename" json:"filename"`
	} `db:"logo" json:"logo"`
	IsVerified         bool                       `db:"is_verified" json:"is_verified"`
	VerificationStatus *KybVerificationStatusType `db:"verification_status" json:"verification_status"`
	UpdatedAt          time.Time                  `db:"updated_at" json:"updated_at"`
	CreatedAt          time.Time                  `db:"created_at" json:"created_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `db:"id" json:"id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	User           *User     `db:"user" json:"user"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (o *Organization) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(o)
}

func (o *Organization) Create(ctx context.Context, userID uuid.UUID) error {
	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}
	rows, err := database.TxQuery(ctx, tx, "organizations/create",
		o.Name, o.Description, o.LogoID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for rows.Next() {
		if err := o.Scan(rows); err != nil {
			tx.Rollback()
			return err
		}
	}
	rows.Close()
	// Creating default member
	rows, err = database.TxQuery(ctx, tx, "organizations/add_member",
		userID, o.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func (o *Organization) NewDID(ctx context.Context) error {
	if o.DID != nil {
		return nil
	}
	did, err := wallet.CreateDID()
	if err != nil {
		return err
	}
	rows, err := database.Query(
		ctx, "organizations/update_did",
		o.ID, did,
	)
	if err != nil {
		return err
	}
	o.DID = &did
	log.Printf("New DID created for `%s` : %s\n", o.Name, *o.DID)
	defer rows.Close()
	return nil
}

func (o *Organization) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx, "organizations/update",
		o.ID, o.Name, o.Description, o.LogoID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := o.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func (o *Organization) UpdateVerification(ctx context.Context, isVerified bool) error {
	rows, err := database.Query(
		ctx, "organizations/update_verification",
		o.ID, isVerified,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := o.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func (*OrganizationMember) TableName() string {
	return "organization_members"
}

func (*OrganizationMember) FetchQuery() string {
	return "organizations/fetch_members"
}

func (m *OrganizationMember) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(m)
}

func GetOrg(id uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id); err != nil {
		return nil, err
	}
	return o, nil
}

func GetOrgByMember(id, userID uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Get(o, "organizations/fetch_one_by_member", id, userID); err != nil {
		return nil, err
	}
	return o, nil
}

func GetOrgsByMember(userID uuid.UUID) ([]Organization, error) {
	var orgs []Organization
	rows, err := database.Queryx("organizations/fetch_by_member", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		o := new(Organization)
		if err := o.Scan(rows); err != nil {
			return nil, err
		}
		orgs = append(orgs, *o)
	}

	return orgs, nil
}
