// Package iface contains all the interface of the system
//go:generate ../../mock.sh
package iface

import (
	"context"
	"database/sql"

	"github.com/rafaelsq/boiler/pkg/entity"
)

// Storage is the storage system
type Storage interface {
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
