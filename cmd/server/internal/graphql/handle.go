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

func ExplorerHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `<!doctype html><html>
<head>
	<style>body {height: 100%; margin: 0; width: 100%; overflow: hidden;}
	#graphiql {height: 100vh;}</style>
	<script crossorigin src="//unpkg.com/react@16/umd/react.development.js"></script>
	<script crossorigin src="//unpkg.com/react-dom@16/umd/react-dom.development.js"></script>
	<link href="//graphiql-test.netlify.app/graphiql.min.css" rel="stylesheet">
	<script src="//graphiql-test.netlify.app/graphiql.min.js"></script>
</head>
<body>
	<div id="graphiql">Loading...</div>
	<script>
	// Parse the search string to get url parameters.
	var search = window.location.search;
	var parameters = {};
	search
	  .substr(1)
	  .split('&')
	  .forEach(function (entry) {
		var eq = entry.indexOf('=');
		if (eq >= 0) {
		  parameters[decodeURIComponent(entry.slice(0, eq))] = decodeURIComponent(
			entry.slice(eq + 1),
		  );
		}
	  });

	// If variables was provided, try to format it.
	if (parameters.variables) {
	  try {
		parameters.variables = JSON.stringify(
		  JSON.parse(parameters.variables),
		  null,
		  2,
		);
	  } catch (e) {
		// Do nothing, we want to display the invalid JSON as a string, rather
		// than present an error.
	  }
	}

	// If headers was provided, try to format it.
	if (parameters.headers) {
	  try {
		parameters.headers = JSON.stringify(
		  JSON.parse(parameters.headers),
		  null,
		  2,
		);
	  } catch (e) {
		// Do nothing, we want to display the invalid JSON as a string, rather
		// than present an error.
	  }
	}

	// When the query and variables string is edited, update the URL bar so
	// that it can be easily shared.
	function onEditQuery(newQuery) {
	  parameters.query = newQuery;
	  updateURL();
	}

	function onEditVariables(newVariables) {
	  parameters.variables = newVariables;
	  updateURL();
	}

	function onEditHeaders(newHeaders) {
	  parameters.headers = newHeaders;
	  updateURL();
	}

	function onEditOperationName(newOperationName) {
	  parameters.operationName = newOperationName;
	  updateURL();
	}

	function updateURL() {
	  var newSearch =
		'?' +
		Object.keys(parameters)
		  .filter(function (key) {
			return Boolean(parameters[key]);
		  })
		  .map(function (key) {
			return (
			  encodeURIComponent(key) + '=' + encodeURIComponent(parameters[key])
			);
		  })
		  .join('&');
	  history.replaceState(null, null, newSearch);
	}

	const isDev = window.location.hostname.match(/localhost$/);
	//const api = isDev ? '/graphql' : '/.netlify/functions/schema-demo';
	let api = '/graphql/query';

	// Render <GraphiQL /> into the body.
	// See the README in the top level of this module to learn more about
	// how you can customize GraphiQL by providing different values or
	// additional child elements.
	ReactDOM.render(
	  React.createElement(GraphiQL, {
		fetcher: GraphiQL.createFetcher({ url: api }),
		query: parameters.query,
		variables: parameters.variables,
		headers: parameters.headers,
		operationName: parameters.operationName,
		onEditQuery: onEditQuery,
		onEditVariables: onEditVariables,
		onEditHeaders: onEditHeaders,
		defaultSecondaryEditorOpen: true,
		onEditOperationName: onEditOperationName,
		headerEditorEnabled: true,
		shouldPersistHeaders: true,
	  }),
	  document.getElementById('graphiql'),
	);
	</script>
</body></html>`)
	}
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
