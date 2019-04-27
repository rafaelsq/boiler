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

func (r *Repo) ByUserID(ctx context.Context, userID uint) ([]*entity.Email, error) {
	return r.db.EmailsByUserID(ctx, userID)
}
