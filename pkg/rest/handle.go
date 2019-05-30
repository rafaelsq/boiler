package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func UsersHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := us.List(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func UserHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err == nil && userID > 0 {
			user, err := us.ByID(r.Context(), int(userID))
			if err == nil {
				if err := json.NewEncoder(w).Encode(user); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

func EmailsHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err == nil && userID > 0 {
			emails, err := es.ByUserID(r.Context(), int(userID))
			if err == nil {
				if err := json.NewEncoder(w).Encode(emails); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}
