package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/handler"
	graphql "github.com/rafaelsq/boiler/pkg/graphql/internal"
	"github.com/rafaelsq/boiler/pkg/iface"
)

// PlayHandle handle Playground 
func PlayHandle() http.HandlerFunc {
	return handler.Playground("Users", "/graphql/query")
}

// QueryHandleFunc return an http HandlerFunc
func QueryHandleFunc(service iface.Service) http.HandlerFunc {
	return handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{
			Resolvers: graphql.NewResolver(service),
		}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("internal server error")
		}),
	)
}
