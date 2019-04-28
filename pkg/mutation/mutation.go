package mutation

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/service"
)

func NewMutation(db entity.DB) *Mutation {
	return &Mutation{
		db:           db,
		emailService: service.NewEmail(db),
	}
}

type Mutation struct {
	db           entity.DB
	emailService service.Email
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.emailService.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: int(input.UserID)}, nil
}
