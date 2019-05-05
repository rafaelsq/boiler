package mutation

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewMutation(service iface.EmailService) *Mutation {
	return &Mutation{
		service: service,
	}
}

type Mutation struct {
	service iface.EmailService
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.service.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: int(input.UserID)}, nil
}
