package mutation

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewMutation(db storage.DB) *Mutation {
	return &Mutation{
		db:           db,
		emailService: service.NewEmail(db),
	}
}

type Mutation struct {
	db           storage.DB
	emailService service.Email
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.emailService.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: int(input.UserID)}, nil
}
