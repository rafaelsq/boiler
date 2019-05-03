package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewUser(repo entity.UserRepository) entity.UserService {
	return &user{
		repo: repo,
	}
}

type user struct {
	repo entity.UserRepository
}

func (s *user) List(ctx context.Context) ([]*entity.User, error) {
	return s.repo.List(ctx)
}

func (s *user) ByID(ctx context.Context, userID int) (*entity.User, error) {
	return s.repo.ByID(ctx, userID)
}

func (s *user) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.ByEmail(ctx, email)
}
