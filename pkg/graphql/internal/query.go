package graphql

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
)

// NewQuery return a new QueryResolver
func NewQuery(ru *resolver.User) QueryResolver {
	return &Query{
		ru: ru,
	}
}

// Query is a Query User struct
type Query struct {
	ru *resolver.User
}

// Users return users
func (r *Query) Users(ctx context.Context, limit *int) ([]*entity.User, error) {
	return r.ru.Users(ctx, uint(*limit))
}

// User return an user
func (r *Query) User(ctx context.Context, userID string) (*entity.User, error) {
	return r.ru.User(ctx, userID)
}
