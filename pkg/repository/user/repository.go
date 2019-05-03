package user

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func New(db storage.DB) entity.UserRepository {
	return &repository{db}
}

type repository struct {
	db storage.DB
}

func (r *repository) List(ctx context.Context) ([]*entity.User, error) {
	return r.db.Users(ctx)
}

func (r *repository) ByID(ctx context.Context, userID int) (*entity.User, error) {
	return r.db.UserByID(ctx, userID)
}

func (r *repository) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.db.UserByEmail(ctx, email)
}
