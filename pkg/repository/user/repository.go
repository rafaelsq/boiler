package user

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func New(storage iface.Storage) iface.UserRepository {
	return &repository{storage}
}

type repository struct {
	storage iface.Storage
}

func (r *repository) List(ctx context.Context) ([]*entity.User, error) {
	return r.storage.Users(ctx)
}

func (r *repository) ByID(ctx context.Context, userID int) (*entity.User, error) {
	return r.storage.UserByID(ctx, userID)
}

func (r *repository) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.storage.UserByEmail(ctx, email)
}
