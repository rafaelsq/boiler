package iface

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type UserRepository interface {
	ByID(context.Context, int) (*entity.User, error)
	ByEmail(context.Context, string) (*entity.User, error)
	List(context.Context) ([]*entity.User, error)
}

type EmailRepository interface {
	ByUserID(context.Context, int) ([]*entity.Email, error)
	Add(context.Context, int, string) (int, error)
}
