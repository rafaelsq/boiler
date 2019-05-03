package resolver

import (
	"context"

	ent "github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
)

func NewEmail(service ent.UserService) *Email {
	return &Email{
		service: service,
	}
}

type Email struct {
	service ent.UserService
}

func (r *Email) ID(ctx context.Context, e *entity.Email) (int, error) {
	return int(e.ID), nil
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	u, err := r.service.ByEmail(ctx, e.Address)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, err
}
