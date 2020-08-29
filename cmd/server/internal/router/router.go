package router

import (
	"compress/flate"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rafaelsq/boiler/cmd/server/internal/rest"
	"github.com/rafaelsq/boiler/cmd/server/internal/website"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/iface"
)

// ApplyMiddlewares add middlewares to the router
func ApplyMiddlewares(r chi.Router, service iface.Service) {
	r.Use(Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(flate.BestCompression))
	r.Use(middleware.Timeout(5 * time.Second))

	// custom middlewares
	r.Use(service.AuthUserMiddleware)

	// Rest Debug
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if len(r.URL.Query()["debug"]) != 0 {
				ctx := r.Context()
				ctx = context.WithValue(ctx, iface.ContextKeyDebug{}, struct{}{})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})
}

// ApplyRoute define the routes of the service
func ApplyRoute(r chi.Router, service iface.Service) {
	// website
	r.Get("/", website.Handle)
	r.Get("/favicon.ico", http.NotFound)
	r.Handle("/static/*", http.FileServer(http.Dir("cmd/server/internal/website")))

	// graphql
	r.Route("/graphql", func(g chi.Router) {
		g.Get("/play", graphql.PlayHandle())
		g.HandleFunc("/query", graphql.QueryHandleFunc(service))
	})

	// rest
	r.Route("/rest", func(r chi.Router) {
		r.Get("/users", rest.ListUsersHandle(service))
		r.Post("/users", rest.AddUserHandle(service))
		r.Get("/users/{userID:[0-9]+}", rest.GetUserHandle(service))
		r.Delete("/users/{userID:[0-9]+}", rest.DeleteUserHandle(service))
		r.Post("/users/login", rest.LoginHandle(service))

		r.Get("/emails", rest.ListEmailsHandle(service))
		r.Post("/emails", rest.AddEmailHandle(service))
		r.Delete("/emails/{emailID:[0-9]+}", rest.DeleteEmailHandle(service))
	})
}
