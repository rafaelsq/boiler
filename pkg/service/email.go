package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewEmail(repo entity.EmailRepository) entity.EmailService {
	return &email{
		repo: repo,
	}
}

type email struct {
	repo entity.EmailRepository
}

func (s *email) Add(ctx context.Context, userID int, address string) (int, error) {
	return s.repo.Add(ctx, userID, address)
}

func (s *email) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return s.repo.ByUserID(ctx, userID)
}
