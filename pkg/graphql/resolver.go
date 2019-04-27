//go:generate go run github.com/99designs/gqlgen
package graphql

import (
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/resolver"
)

func NewResolver(db entity.DB) *Resolver {
	return &Resolver{db}
}

type Resolver struct {
	db entity.DB
}

func (r *Resolver) Query() QueryResolver {
	return NewQuery(r.db)
}

func (r *Resolver) Mutation() MutationResolver {
	return NewMutation(r.db)
}

func (r *Resolver) User() UserResolver {
	return resolver.NewUser(r.db)
}

func (r *Resolver) Email() EmailResolver {
	return resolver.NewEmail(r.db)
}
