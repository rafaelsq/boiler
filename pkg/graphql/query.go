package graphql

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/resolver"
)

func NewQuery(db entity.DB) *Query {
	return &Query{
		db:           db,
		userResolver: resolver.NewUser(db),
	}
}

type Query struct {
	db           entity.DB
	userResolver *resolver.User
}

func (r *Query) Users(ctx context.Context) ([]*entity.User, error) {
	return r.userResolver.Users(ctx)
}

func (r *Query) User(ctx context.Context, userID int) (*entity.User, error) {
	return r.userResolver.User(ctx, userID)
}
