package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/service"
)

func NewEmail(db entity.DB) *Email {
	return &Email{
		db:          db,
		userService: service.NewUser(db),
	}
}

type Email struct {
	db          entity.DB
	userService service.User
}

func (r *Email) ID(ctx context.Context, e *entity.Email) (int, error) {
	return int(e.ID), nil
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	return r.userService.ByEmail(ctx, e.Address)
}
