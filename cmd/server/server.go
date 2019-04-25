package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase"
)

var port = flag.Int("port", 2000, "")

func main() {
	flag.Parse()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()[1:]

		if path == "" {
			fmt.Fprintf(w, "<h1>Home</h1><a href=\"/users\">list users</a>")
			return
		}

		if userID, err := strconv.Atoi(path); err == nil && userID > 0 {
			ucase := usecase.NewUser( /*db*/ )
			if user, err := ucase.ByID(r.Context(), userID); err == nil && user != nil {
				fmt.Fprintf(w, "<h1>User</h1><p>%d - %s", user.ID, user.Name)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>not found</h1><a href=\"/\">home</a>")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		ucase := usecase.NewUser( /*db*/ )
		users, err := ucase.Filter(r.Context(), &entity.UserFilter{
			Order: entity.UserFilterOrderASC,
		})
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
