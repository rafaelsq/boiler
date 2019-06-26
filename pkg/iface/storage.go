//go:generate ../../mock.sh
package iface

import (
	"context"
	"database/sql"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type Storage interface {
	// begin transaction
	Tx() (*sql.Tx, error)

	// user
	AddUser(ctx context.Context, tx *sql.Tx, name string) (int, error)
	DeleteUser(ctx context.Context, userID int) error
	FilterUsers(ctx context.Context, filter FilterUsers) ([]*entity.User, error)

	// email
	AddEmail(ctx context.Context, tx *sql.Tx, userID int, address string) (int, error)
	DeleteEmail(ctx context.Context, emailID int) error
	FilterEmails(ctx context.Context, filter FilterEmails) ([]*entity.Email, error)
}
