//go:generate mockgen -package=mock -source=$GOFILE -destination=mock/store.go
package store

import (
	"context"
	"database/sql"
	"errors"

	"boiler/pkg/entity"
)

var (
	// ErrNotFound not found error
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists already exists error
	ErrAlreadyExists = errors.New("already exists")
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
	AddUser(ctx context.Context, tx *sql.Tx, name, password string) (int64, error)
	DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error
	FilterUsersID(ctx context.Context, filter FilterUsers) ([]int64, error)
	FetchUsers(ctx context.Context, ID ...int64) ([]*entity.User, error)

	// email
	AddEmail(ctx context.Context, tx *sql.Tx, userID int64, address string) (int64, error)
	DeleteEmail(ctx context.Context, tx *sql.Tx, emailID int64) error
	DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error
	FilterEmails(ctx context.Context, filter FilterEmails) ([]*entity.Email, error)
}
