//go:generate ../../mock.sh service
package iface

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type UserService interface {
	ByID(context.Context, int) (*entity.User, error)
	ByEmail(context.Context, string) (*entity.User, error)
	List(context.Context) ([]*entity.User, error)
}

type EmailService interface {
	ByUserID(context.Context, int) ([]*entity.Email, error)
	Add(context.Context, int, string) (int, error)
}
