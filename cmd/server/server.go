package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/99designs/gqlgen/handler"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/graphql/resolver"
	"github.com/rafaelsq/boiler/pkg/storage"
)

var port = flag.Int("port", 2000, "")

func main() {
	flag.Parse()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.Handle("/play", handler.Playground("Users", "/query"))

	http.HandleFunc("/query", handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{
			Resolvers: graphql.NewResolver(storage.GetDB()),
		}),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			log.Print(err)
			debug.PrintStack()
			return errors.New("internal server error")
		}),
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()[1:]

		if path == "" {
			fmt.Fprintf(w, "<h1>Home</h1><p><a href=\"/users\">list users</a></p><p><a href=\"/play\">GraphQL</a></p>")
			return
		}

		if userID, err := strconv.ParseUint(path, 10, 64); err == nil && userID > 0 {
			ucase := resolver.NewUser(storage.GetDB())
			if user, err := ucase.User(r.Context(), int(userID)); err == nil && user != nil {
				fmt.Fprintf(w, "<h1>User</h1><p>%d - %s</p><ul>", user.ID, user.Name)
				for _, email := range user.Emails {
					fmt.Fprintf(w, "<li>%d - <%s>%s</li>", email.ID, email.User.Name, email.Address)
				}
				fmt.Fprintf(w, "</ul>")
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>not found</h1><a href=\"/\">home</a>")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		ucase := resolver.NewUser(storage.GetDB())
		users, err := ucase.Users(r.Context())
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

	log.Printf("Listening on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
