package email

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func New(db entity.DB) Repository {
	return &repository{db}
}

type Repository interface {
	ByUserID(context.Context, int) ([]*entity.Email, error)
	Add(context.Context, int, string) (int, error)
}

type repository struct {
	db entity.DB
}

func (r *repository) Add(ctx context.Context, userID int, address string) (int, error) {
	return r.db.AddEmail(ctx, userID, address)
}

func (r *repository) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return r.db.EmailsByUserID(ctx, userID)
}
