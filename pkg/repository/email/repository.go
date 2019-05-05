package email

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func New(storage iface.Storage) iface.EmailRepository {
	return &repository{storage}
}

type repository struct {
	storage iface.Storage
}

func (r *repository) Add(ctx context.Context, userID int, address string) (int, error) {
	return r.storage.AddEmail(ctx, userID, address)
}

func (r *repository) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return r.storage.EmailsByUserID(ctx, userID)
}
