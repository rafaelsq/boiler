package mutation

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/repository/email"
)

func NewMutation(db entity.DB) *Mutation {
	return &Mutation{
		db:        db,
		emailRepo: email.NewRepo(db),
	}
}

type Mutation struct {
	db        entity.DB
	emailRepo entity.EmailRepository
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.emailRepo.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: int(input.UserID)}, nil
}
