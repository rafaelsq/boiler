package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/handler"
	graphql "github.com/rafaelsq/boiler/pkg/graphql/internal"
	er "github.com/rafaelsq/boiler/pkg/repository/email"
	ur "github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func NewPlayHandle() http.Handler {
	return handler.Playground("Users", "/query")
}

func NewHandleFunc(db storage.DB) func(http.ResponseWriter, *http.Request) {
	us := service.NewUser(ur.New(db))
	es := service.NewEmail(er.New(db))

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
