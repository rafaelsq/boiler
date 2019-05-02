package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/entity"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewEmail(db storage.DB) *Email {
	return &Email{
		db:          db,
		userService: service.NewUser(db),
	}
}

type Email struct {
	db          storage.DB
	userService service.User
}

func (r *Email) ID(ctx context.Context, e *entity.Email) (int, error) {
	return int(e.ID), nil
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	u, err := r.userService.ByEmail(ctx, e.Address)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, err
}
