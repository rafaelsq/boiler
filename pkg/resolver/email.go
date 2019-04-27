package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/resolver/internal/repository/email"
	"github.com/rafaelsq/boiler/pkg/resolver/internal/repository/user"
)

func NewEmail(db entity.DB) *Email {
	return &Email{
		db:        db,
		userRepo:  user.NewRepo(db),
		emailRepo: email.NewRepo(db),
	}
}

type Email struct {
	db        entity.DB
	userRepo  entity.UserRepository
	emailRepo entity.EmailRepository
}

func (r *Email) ID(ctx context.Context, e *entity.Email) (int, error) {
	return int(e.ID), nil
}

func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	return r.userRepo.ByEmail(ctx, e.Address)
}
