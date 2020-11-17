package router

import (
	"compress/flate"
	"context"
	"net/http"
	"time"

	"boiler/cmd/server/internal/graphql"
	"boiler/cmd/server/internal/rest"
	"boiler/cmd/server/internal/website"
	"boiler/pkg/iface"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// ApplyMiddlewares add middlewares to the router
func ApplyMiddlewares(r chi.Router, service iface.Service) {
	r.Use(Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
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
		g.Handle("/query", graphql.QueryHandler(service))
	})

	// rest
	r.Route("/rest", func(r chi.Router) {
		h := rest.New(service)

		r.Get("/users", h.ListUsers)
		r.Post("/users", h.AddUser)
		r.Get("/users/{userID:[0-9]+}", h.GetUser)
		r.Delete("/users/{userID:[0-9]+}", h.DeleteUser)
		r.Post("/users/login", h.AuthUser)

		r.Get("/emails", h.ListEmails)
		r.Post("/emails", h.AddEmail)
		r.Delete("/emails/{emailID:[0-9]+}", h.DeleteEmail)
	})
}
