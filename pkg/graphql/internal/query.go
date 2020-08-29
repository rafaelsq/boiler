package graphql

import (
	"context"
	"errors"
	"strconv"

	globalEntity "github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/iface"
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

func (r *Query) Viewer(ctx context.Context) (*entity.User, error) {

	raw := ctx.Value(iface.ContextKeyAuthenticationUser{})
	if raw == nil {
		return nil, errors.New("unauthorized")
	}

	return r.ru.User(ctx, strconv.FormatInt(raw.(*globalEntity.JWTUser).ID, 10))
}
