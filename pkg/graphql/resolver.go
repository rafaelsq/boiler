//go:generate go run github.com/99designs/gqlgen
package graphql

import (
	"github.com/rafaelsq/boiler/pkg/graphql/mutation"
	"github.com/rafaelsq/boiler/pkg/graphql/resolver"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewResolver(db storage.DB) ResolverRoot {
	return &Resolver{db}
}

type Resolver struct {
	db storage.DB
}

func (r *Resolver) Query() QueryResolver {
	return NewQuery(r.db)
}

func (r *Resolver) Mutation() MutationResolver {
	return mutation.NewMutation(r.db)
}

func (r *Resolver) User() UserResolver {
	return resolver.NewUser(r.db)
}

func (r *Resolver) Email() EmailResolver {
	return resolver.NewEmail(r.db)
}
