package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	er "github.com/rafaelsq/boiler/pkg/service/internal/repository/email"
)

func NewEmail(db entity.DB) Email {
	return &email{
		db:        db,
		emailRepo: er.New(db),
	}
}

type Email interface {
	ByUserID(context.Context, int) ([]*entity.Email, error)
	Add(context.Context, int, string) (int, error)
}

type email struct {
	db        entity.DB
	emailRepo er.Repository
}

func (s *email) Add(ctx context.Context, userID int, address string) (int, error) {
	return s.emailRepo.Add(ctx, userID, address)
}

func (s *email) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return s.emailRepo.ByUserID(ctx, userID)
}
