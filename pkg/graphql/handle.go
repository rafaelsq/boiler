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
	er "github.com/rafaelsq/boiler/pkg/repository/email"
	ur "github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/rafaelsq/boiler/pkg/service"
)

func NewPlayHandle() http.HandlerFunc {
	return handler.Playground("Users", "/query")
}

func NewHandleFunc(storage iface.Storage) http.HandlerFunc {
	us := service.NewUser(ur.New(storage))
	es := service.NewEmail(er.New(storage))

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
