package user

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewRepo(db entity.DB) entity.UserRepository {
	return &Repo{db}
}

type Repo struct {
	db entity.DB
}

func (r *Repo) Users(ctx context.Context) ([]*entity.User, error) {
	return r.db.Users(ctx)
}

func (r *Repo) ByID(ctx context.Context, userID uint) (*entity.User, error) {
	return r.db.UserByID(ctx, userID)
}

func (r *Repo) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.db.UserByEmail(ctx, email)
}
