package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"boiler/pkg/iface"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

// PlayHandle handle Playground
func PlayHandle() http.HandlerFunc {
	return playground.Handler("Users", "/graphql/query")
}

// QueryHandleFunc return an http HandlerFunc
func QueryHandler(service iface.Service) http.Handler {
	hldr := handler.New(
		NewExecutableSchema(Config{
			Resolvers: NewResolver(service),
		}),
	)

	hldr.AddTransport(transport.GET{})
	hldr.AddTransport(transport.POST{})

	hldr.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Print(err)
		debug.PrintStack()
		return errors.New("internal server error")
	})

	return hldr
}
