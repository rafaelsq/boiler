package main

import (
	"compress/flate"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/rest"
	"github.com/rafaelsq/boiler/pkg/storage"
	"github.com/rafaelsq/boiler/pkg/website"
)

var port = flag.Int("port", 2000, "")

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Compress(flate.BestCompression))
	r.Use(middleware.Timeout(60 * time.Second))

	// injections
	r.Use(middleware.WithValue("storage", storage.GetDB()))

	// graphql
	r.Route("/graphql", func(r chi.Router) {
		r.Get("/play", graphql.NewPlayHandle())
		r.HandleFunc("/query", graphql.NewHandleFunc(storage.GetDB()))
	})

	// website
	r.Get("/", website.Handle)
	r.Get("/favicon.ico", http.NotFound)
	r.Handle("/static/*", http.FileServer(http.Dir("pkg/website")))

	// rest
	r.Route("/rest", func(r chi.Router) {
		r.Get("/users", rest.UsersHandle)
		r.Get("/user/{userID:[0-9]+}", rest.UserHandle)
		r.Get("/emails/{userID:[0-9]+}", rest.EmailsHandle)

	})

	// gracefull shutdown
	srv := http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: r}

	c := make(chan os.Signal, 1)
	iddleConnections := make(chan struct{})
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down..")

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// start http shutdown
			if err := srv.Shutdown(ctx); err != nil {
				log.Println("shutdown error", err)
			}

			close(iddleConnections)
		}
	}()

	log.Printf("Listening on :%d\n", *port)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	log.Println("waiting iddle connections...")
	<-iddleConnections
	log.Println("done, bye!")
}
