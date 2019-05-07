//go:generate ../../mock.sh storage
package iface

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type Storage interface {
	Users(context.Context) ([]*entity.User, error)
	UserByID(context.Context, int) (*entity.User, error)
	UserByEmail(context.Context, string) (*entity.User, error)
	AddEmail(context.Context, int, string) (int, error)
	EmailsByUserID(context.Context, int) ([]*entity.Email, error)
}
