package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"shin/src/wallet"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type Organization struct {
	ID                 uuid.UUID                  `db:"id" json:"id"`
	DID                *string                    `db:"did" json:"did"`
	Name               string                     `db:"name" json:"name"`
	Description        string                     `db:"description" json:"description"`
	Logo               *Media                     `db:"-" json:"logo"`
	LogoJson           types.JSONText             `db:"logo" json:"-"`
	Verified           bool                       `db:"verified" json:"verified"`
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

	if o.Logo != nil {
		b, _ := json.Marshal(o.Logo)
		o.LogoJson.Scan(b)
	}

	if o.ID == uuid.Nil {
		newID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		o.ID = newID
	}

	rows, err := database.TxQuery(
		ctx,
		tx,
		"organizations/create",
		o.ID,
		o.Name,
		o.Description,
		o.LogoJson,
		o.Verified,
	)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(o); err != nil {
			tx.Rollback()
			return err
		}
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

	if o.Logo != nil {
		b, _ := json.Marshal(o.Logo)
		o.LogoJson.Scan(b)
	}

	rows, err := database.Query(
		ctx, "organizations/update",
		o.ID, o.Name, o.Description, o.LogoJson,
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
	var (
		orgs      = []Organization{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("organizations/fetch_by_member", &fetchList, userID); err != nil {
		return orgs, err
	}

	if len(fetchList) < 1 {
		return orgs, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&orgs, ids...); err != nil {
		return orgs, err
	}
	return orgs, nil
}
