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
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/iface"
	er "github.com/rafaelsq/boiler/pkg/repository/email"
	ur "github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
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

	r.Get("/favicon.ico", http.NotFound)

	// graphql
	r.Get("/play", graphql.NewPlayHandle())
	r.HandleFunc("/query", graphql.NewHandleFunc(storage.GetDB()))

	// html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Home</h1><p><a href=\"/users\">list users</a></p><p><a href=\"/play\">GraphQL</a></p>")
	})

	r.Get("/{userID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		st := r.Context().Value("storage").(iface.Storage)
		if userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err == nil && userID > 0 {
			ucase := service.NewUser(ur.New(st))
			cEmails := make(chan []*entity.Email)
			go func() {
				es, _ := service.NewEmail(er.New(st)).ByUserID(r.Context(), int(userID))
				if err == nil {
					cEmails <- es
					return
				}
				panic(err)
			}()

			if user, err := ucase.ByID(r.Context(), int(userID)); err == nil && user != nil {
				fmt.Fprintf(w, "<h1>User</h1><p>%d - %s</p><ul>", user.ID, user.Name)
				for _, email := range <-cEmails {
					fmt.Fprintf(w, "<li>%d - <%s>%s</li>", email.ID, user.Name, email.Address)
				}
				fmt.Fprintf(w, "</ul>")
				return
			}
		}
	})

	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(time.Second * 4):
			fmt.Println("by time")
		case <-r.Context().Done():
			fmt.Println("timeout")
		}
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		ucase := service.NewUser(ur.New(r.Context().Value("storage").(iface.Storage)))
		users, err := ucase.List(r.Context())
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "<h1>Users</h1><ul>")
		for _, user := range users {
			fmt.Fprintf(w, "<li><a href=\"/%d\">%s<a/></li>", user.ID, user.Name)
		}
		fmt.Fprintf(w, "</ul>")
	})

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
