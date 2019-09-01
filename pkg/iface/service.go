//go:generate ../../mock.sh
package iface

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type Service interface {
	// user
	AddUser(context.Context, string) (int64, error)
	DeleteUser(context.Context, int64) error
	FilterUsers(context.Context, FilterUsers) ([]*entity.User, error)
	GetUserByID(context.Context, int64) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)

	// email
	FilterEmails(context.Context, FilterEmails) ([]*entity.Email, error)
	AddEmail(context.Context, int64, string) (int64, error)
	DeleteEmail(context.Context, int64) error
}
