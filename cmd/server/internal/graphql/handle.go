package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"boiler/pkg/iface"

	"github.com/99designs/gqlgen/handler"
)

// PlayHandle handle Playground
func PlayHandle() http.HandlerFunc {
	return handler.Playground("Users", "/graphql/query")
}

// QueryHandleFunc return an http HandlerFunc
func QueryHandleFunc(service iface.Service) http.HandlerFunc {
	return handler.GraphQL(
		NewExecutableSchema(Config{
			Resolvers: NewResolver(service),
		}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("internal server error")
		}),
	)
}
