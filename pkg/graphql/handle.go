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

func PlayHandle() http.HandlerFunc {
	return handler.Playground("Users", "/graphql/query")
}

func QueryHandleFunc(us iface.UserService, es iface.EmailService) http.HandlerFunc {
	return handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{
			Resolvers: graphql.NewResolver(us, es),
		}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("internal server error")
		}),
	)
}
