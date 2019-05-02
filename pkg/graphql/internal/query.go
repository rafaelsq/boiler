package graphql

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewQuery(db storage.DB) QueryResolver {
	return &Query{
		db:           db,
		userResolver: resolver.NewUser(db),
	}
}

type Query struct {
	db           storage.DB
	userResolver *resolver.User
}

func (r *Query) Users(ctx context.Context) ([]*entity.User, error) {
	return r.userResolver.Users(ctx)
}

func (r *Query) User(ctx context.Context, userID int) (*entity.User, error) {
	return r.userResolver.User(ctx, userID)
}
