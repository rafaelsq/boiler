package mutation

import (
	"context"

	ent "github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
)

func NewMutation(service ent.EmailService) *Mutation {
	return &Mutation{
		service: service,
	}
}

type Mutation struct {
	service ent.EmailService
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.service.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: int(input.UserID)}, nil
}
