package graphql

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"

	"boiler/pkg/iface"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog/log"
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
		log.Error().Interface("err", err).Caller().Send()
		debug.PrintStack()
		return errors.New("internal server error")
	})

	return hldr
}
