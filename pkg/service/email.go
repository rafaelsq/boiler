package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewEmail(repo iface.EmailRepository) iface.EmailService {
	return &email{
		repo: repo,
	}
}

type email struct {
	repo iface.EmailRepository
}

func (s *email) Add(ctx context.Context, userID int, address string) (int, error) {
	return s.repo.Add(ctx, userID, address)
}

func (s *email) Delete(ctx context.Context, emailID int) error {
	return s.repo.Delete(ctx, emailID)
}

func (s *email) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return s.repo.ByUserID(ctx, userID)
}
