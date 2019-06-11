package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func AddUserHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Name string `json:"name"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not parse payload")
			return
		}

		payload.Name = strings.TrimSpace(payload.Name)
		if len(payload.Name) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "empty name")
			return
		}

		userID, err := us.Add(r.Context(), payload.Name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "service fail")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct{ UserID int }{userID}); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not encode response")
		}
	}
}

func UsersHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		rawLimit := r.URL.Query()["limit"]
		if len(rawLimit) > 0 {
			var err error
			limit, err = strconv.Atoi(rawLimit[0])
			if err != nil || limit <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "invalid limit \"%s\"", rawLimit[0])
				return
			}
		}

		users, err := us.List(r.Context(), uint(limit))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "service fail")
			return
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not encode response")
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

func AddEmailHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			UserID  int    `json:"user"`
			Address string `json:"address"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload.Address = strings.TrimSpace(payload.Address)
		if len(payload.Address) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "empty address")
			return
		}
		if payload.UserID < 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid user ID")
			return
		}

		emailID, err := es.Add(r.Context(), payload.UserID, payload.Address)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct{ EmailID int }{emailID}); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
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
