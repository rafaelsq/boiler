package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
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
			Fail(w, r, http.StatusBadRequest, "could not parse payload")
			return
		}

		payload.Name = strings.TrimSpace(payload.Name)
		if len(payload.Name) == 0 {
			Fail(w, r, http.StatusBadRequest, "empty name")
			return
		}

		userID, err := us.Add(r.Context(), payload.Name)
		if err != nil {
			log.Println(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, struct{ UserID int }{userID})
	}
}

func ListUsersHandle(us iface.UserService) http.HandlerFunc {
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
			fmt.Fprintf(w, "service failed")
			return
		}

		JSON(w, r, users)
	}
}

func GetUserHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		user, err := us.ByID(r.Context(), int(userID))
		if err != nil {
			log.Println(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, user)
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
			Fail(w, r, http.StatusBadRequest, "invalid payload")
			return
		}

		email, err := mail.ParseAddress(payload.Address)
		if err != nil {
			fmt.Println("...", err)
			Fail(w, r, http.StatusBadRequest, "invalid email address")
			return
		}

		if payload.UserID < 1 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		emailID, err := es.Add(r.Context(), payload.UserID, email.Address)
		if err != nil {
			log.Println(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, struct {
			EmailID int `json:"emailID"`
		}{emailID})
	}
}
func ListUserEmailsHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid user id")
			return
		}

		emails, err := es.ByUserID(r.Context(), int(userID))
		if err != nil {
			log.Println(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, emails)
	}
}
