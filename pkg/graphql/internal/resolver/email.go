package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewEmail(service iface.Service) *Email {
	return &Email{
		service: service,
	}
}

type Email struct {
	service iface.Service
}

func (r *Email) ID(ctx context.Context, e *entity.Email) (int, error) {
	return e.ID, nil
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	u, err := r.service.GetUserByEmail(ctx, e.Address)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, err
}
