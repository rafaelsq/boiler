package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	ur "github.com/rafaelsq/boiler/pkg/service/internal/repository/user"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewUser(db storage.DB) User {
	return &user{
		db:       db,
		userRepo: ur.New(db),
	}
}

type User interface {
	ByID(context.Context, int) (*entity.User, error)
	ByEmail(context.Context, string) (*entity.User, error)
	List(context.Context) ([]*entity.User, error)
}

type user struct {
	db       storage.DB
	userRepo ur.Repository
}

func (s *user) List(ctx context.Context) ([]*entity.User, error) {
	return s.userRepo.List(ctx)
}

func (s *user) ByID(ctx context.Context, userID int) (*entity.User, error) {
	return s.userRepo.ByID(ctx, userID)
}

func (s *user) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.ByEmail(ctx, email)
}
