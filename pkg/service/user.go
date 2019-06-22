package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewUser(repo iface.UserRepository) iface.UserService {
	return &user{
		repo: repo,
	}
}

type user struct {
	repo iface.UserRepository
}

func (s *user) Add(ctx context.Context, name string) (int, error) {
	return s.repo.Add(ctx, name)
}

func (s *user) Delete(ctx context.Context, userID int) error {
	return s.repo.Delete(ctx, userID)
}

func (s *user) List(ctx context.Context, limit uint) ([]*entity.User, error) {
	return s.repo.List(ctx, limit)
}

func (s *user) ByID(ctx context.Context, userID int) (*entity.User, error) {
	return s.repo.ByID(ctx, userID)
}

func (s *user) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.ByEmail(ctx, email)
}
