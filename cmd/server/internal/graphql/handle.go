package graphql

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"boiler/pkg/errors"
	"boiler/pkg/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// PlayHandle handle Playground
func PlayHandle() http.HandlerFunc {
	return playground.Handler("Play", "/graphql/query")
}

// QueryHandleFunc return an http HandlerFunc
func QueryHandler(service service.Interface) http.Handler {
	hldr := handler.New(
		NewExecutableSchema(Config{
			Resolvers: NewResolver(service),
		}),
	)

	hldr.Use(extension.Introspection{})
	hldr.Use(apollotracing.Tracer{})

	hldr.AddTransport(transport.Options{})
	hldr.AddTransport(transport.GET{})
	hldr.AddTransport(transport.POST{})
	hldr.AddTransport(transport.MultipartForm{})

	hldr.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		if false {
			stackTitle, stackTrace := errors.GetStack()
			fmt.Fprintf(
				os.Stderr,
				"ERROR: %s\n%s\n    %s\n",
				err,
				stackTitle,
				strings.Join(stackTrace, "\n    "),
			)
			return fmt.Errorf("%#v", err)
		}
		log.Error().Str("file", errors.Caller()).Interface("err", err).Send()
		return errors.New("service unavailable")
	})

	hldr.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {

		err := graphql.DefaultErrorPresenter(ctx, e)
		codes := []string{}

		var ce *errors.CodeErr
		errs := e
		for {
			if errors.As(errs, &ce) {
				codes = append(codes, ce.Code)
				errs = errors.Unwrap(ce)
				continue
			}

			break
		}

		if len(codes) == 0 {
			log.Error().Str("file", errors.Caller()).Err(errors.Unwrap(e)).Send()
			codes = append(codes, "INTERNAL_SERVER_ERROR")
			err.Message = "service unavailable"
		}

		err.Extensions = map[string]interface{}{
			"codes": codes,
		}

		return e.(*gqlerror.Error)
	})

	return hldr
}
