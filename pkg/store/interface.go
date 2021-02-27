//go:generate go run github.com/golang/mock/mockgen -package=mock -source=$GOFILE -destination=mock/store.go
package store

import (
	"context"
	"database/sql"

	"boiler/pkg/entity"
)

// FilterUsers is the input for filter users
type FilterUsers struct {
	Email  string
	Offset uint
	Limit  uint
}

// FilterEmails is the input for filter emails
type FilterEmails struct {
	EmailID int64
	UserID  int64
	Offset  uint
	Limit   uint
}

// Interface
type Interface interface {
	// begin transaction
	Tx() (*sql.Tx, error)

	// user
	AddUser(ctx context.Context, tx *sql.Tx, user *entity.User) error
	DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error
	FilterUsersID(ctx context.Context, filter FilterUsers, IDs *[]int64) error
	FetchUsers(ctx context.Context, ID []int64, users *[]entity.User) error

	// email
	AddEmail(ctx context.Context, tx *sql.Tx, email *entity.Email) error
	DeleteEmail(ctx context.Context, tx *sql.Tx, email int64) error
	DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error
	FilterEmails(ctx context.Context, filter FilterEmails, emails *[]entity.Email) error
}
