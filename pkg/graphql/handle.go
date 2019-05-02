package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/handler"
	graphql "github.com/rafaelsq/boiler/pkg/graphql/internal"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewPlayHandle() http.Handler {
	return handler.Playground("Users", "/query")
}

func NewHandleFunc(db storage.DB) func(http.ResponseWriter, *http.Request) {
	return handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{
			Resolvers: graphql.NewResolver(db),
		}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("internal server error")
		}),
	)
}
