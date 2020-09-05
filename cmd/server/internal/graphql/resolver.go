// Package graphql contains all the graphql resources
package graphql

import (
	"boiler/cmd/server/internal/graphql/mutation"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/iface"
)

// NewResolver return a new Resolver
func NewResolver(service iface.Service) ResolverRoot {
	return &Resolver{service}
}

// Resolver is the Resolver of the service
type Resolver struct {
	service iface.Service
}

// Query return a new QueryResolver
func (r *Resolver) Query() QueryResolver {
	return NewQuery(resolver.NewUser(r.service))
}

// Mutation return a new MutationResolver
func (r *Resolver) Mutation() MutationResolver {
	return mutation.NewMutation(r.service)
}

// User return a new UserResolver
func (r *Resolver) User() UserResolver {
	return resolver.NewUser(r.service)
}

// UserResponse return a new UserResponseResolver
func (r *Resolver) UserResponse() UserResponseResolver {
	return resolver.NewResponse(r.service)
}

// AuthUserResponse return a new UserResponseResolver
func (r *Resolver) AuthUserResponse() AuthUserResponseResolver {
	return resolver.NewAuthUserResponse(r.service)
}

// Email return a new EmailResolver
func (r *Resolver) Email() EmailResolver {
	return resolver.NewEmail(r.service)
}

// EmailResponse return a new EmailResponseResolver
func (r *Resolver) EmailResponse() EmailResponseResolver {
	return resolver.NewResponse(r.service)
}
