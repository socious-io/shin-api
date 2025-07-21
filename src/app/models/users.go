package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	Username        string         `db:"username" json:"username"`
	Email           string         `db:"email" json:"email"`
	Password        *string        `db:"password" json:"-"`
	JobTitle        *string        `db:"job_title" json:"job_title"`
	Bio             *string        `db:"bio" json:"bio"`
	FirstName       *string        `db:"first_name" json:"first_name"`
	LastName        *string        `db:"last_name" json:"last_name"`
	Phone           *string        `db:"phone" json:"phone"`
	Avatar          *Media         `db:"-" json:"avatar"`
	AvatarJson      types.JSONText `db:"avatar" json:"-"`
	Status          string         `db:"status" json:"status"`
	PasswordExpired bool           `db:"password_expired" json:"password_expired"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (User) FetchQuery() string {
	return "users/fetch"
}

func (u *User) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/register",
		u.FirstName, u.LastName, u.Username, u.Email, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Upsert(ctx context.Context) error {
	if u.Avatar != nil {
		b, _ := json.Marshal(u.Avatar)
		u.AvatarJson.Scan(b)
	}

	if u.ID == uuid.Nil {
		newID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		u.ID = newID
	}

	rows, err := database.Query(
		ctx,
		"users/upsert",
		u.ID,
		u.FirstName,
		u.LastName,
		u.Username,
		u.Email,
		u.Bio,
		u.Phone,
		u.AvatarJson,
	)

	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Verify(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/verify",
		u.ID, u.Status,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) ExpirePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/expire_password",
		u.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdatePassword(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_password",
		u.ID, u.Password,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) UpdateProfile(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"users/update_profile",
		u.ID, u.FirstName, u.LastName, u.Bio, u.JobTitle, u.Phone, u.Username, u.Avatar,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return database.Fetch(u, u.ID)
}

func GetUser(id uuid.UUID) (*User, error) {
	u := new(User)
	if err := database.Fetch(u, id.String()); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_email", email); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByUsername(username string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_username", username); err != nil {
		return nil, err
	}
	return u, nil
}
