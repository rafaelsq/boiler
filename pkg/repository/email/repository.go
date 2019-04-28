package email

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewRepo(db entity.DB) entity.EmailRepository {
	return &Repo{db}
}

type Repo struct {
	db entity.DB
}

func (r *Repo) Add(ctx context.Context, userID int, address string) (int, error) {
	return r.db.AddEmail(ctx, userID, address)
}

func (r *Repo) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	return r.db.EmailsByUserID(ctx, userID)
}
