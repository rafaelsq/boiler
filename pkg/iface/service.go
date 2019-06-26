//go:generate ../../mock.sh
package iface

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type Service interface {
	// user
	AddUser(context.Context, string) (int, error)
	DeleteUser(context.Context, int) error
	FilterUsers(context.Context, FilterUsers) ([]*entity.User, error)
	GetUserByID(context.Context, int) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)

	// email
	FilterEmails(context.Context, FilterEmails) ([]*entity.Email, error)
	AddEmail(context.Context, int, string) (int, error)
	DeleteEmail(context.Context, int) error
}
