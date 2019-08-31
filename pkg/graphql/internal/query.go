package graphql

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
)

func NewQuery(ru *resolver.User) QueryResolver {
	return &Query{
		ru: ru,
	}
}

type Query struct {
	ru *resolver.User
}

func (r *Query) Users(ctx context.Context, limit *int) ([]*entity.User, error) {
	return r.ru.Users(ctx, uint(*limit))
}

func (r *Query) User(ctx context.Context, userID string) (*entity.User, error) {
	return r.ru.User(ctx, userID)
}
