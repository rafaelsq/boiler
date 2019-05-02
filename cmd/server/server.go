package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/graphql"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

var port = flag.Int("port", 2000, "")

func main() {
	flag.Parse()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.Handle("/play", graphql.NewPlayHandle())
	http.HandleFunc("/query", graphql.NewHandleFunc(storage.GetDB()))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()[1:]

		if path == "" {
			fmt.Fprintf(w, "<h1>Home</h1><p><a href=\"/users\">list users</a></p><p><a href=\"/play\">GraphQL</a></p>")
			return
		}

		if userID, err := strconv.ParseUint(path, 10, 64); err == nil && userID > 0 {
			ucase := service.NewUser(storage.GetDB())
			cEmails := make(chan []*entity.Email)
			go func() {
				es, _ := service.NewEmail(storage.GetDB()).ByUserID(r.Context(), int(userID))
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

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>not found</h1><a href=\"/\">home</a>")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		ucase := service.NewUser(storage.GetDB())
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

	log.Printf("Listening on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
