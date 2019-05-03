//go:generate go run github.com/99designs/gqlgen
package graphql

import (
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/mutation"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
)

func NewResolver(us entity.UserService, es entity.EmailService) ResolverRoot {
	return &Resolver{us, es}
}

type Resolver struct {
	us entity.UserService
	es entity.EmailService
}

func (r *Resolver) Query() QueryResolver {
	return NewQuery(resolver.NewUser(r.us, r.es))
}

func (r *Resolver) Mutation() MutationResolver {
	return mutation.NewMutation(r.es)
}

func (r *Resolver) User() UserResolver {
	return resolver.NewUser(r.us, r.es)
}

func (r *Resolver) Email() EmailResolver {
	return resolver.NewEmail(r.us)
}
