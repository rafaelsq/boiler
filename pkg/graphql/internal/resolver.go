//go:generate go run github.com/99designs/gqlgen
package graphql

import (
	"github.com/rafaelsq/boiler/pkg/graphql/internal/mutation"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewResolver(service iface.Service) ResolverRoot {
	return &Resolver{service}
}

type Resolver struct {
	service iface.Service
}

func (r *Resolver) Query() QueryResolver {
	return NewQuery(resolver.NewUser(r.service))
}

func (r *Resolver) Mutation() MutationResolver {
	return mutation.NewMutation(r.service)
}

func (r *Resolver) User() UserResolver {
	return resolver.NewUser(r.service)
}

func (r *Resolver) UserResponse() UserResponseResolver {
	return resolver.NewResponse(r.service)
}

func (r *Resolver) Email() EmailResolver {
	return resolver.NewEmail(r.service)
}

func (r *Resolver) EmailResponse() EmailResponseResolver {
	return resolver.NewResponse(r.service)
}
