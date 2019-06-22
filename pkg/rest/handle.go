package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rafaelsq/boiler/pkg/errors"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func AddUserHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Name string `json:"name"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
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
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"userID": userID,
		})
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
			errors.Log(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"users": users,
		})
	}
}

func DeleteUserHandle(us iface.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		err = us.Delete(r.Context(), int(userID))
		if err != nil {
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, nil)
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
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"user": user,
		})
	}
}

func AddEmailHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			UserID  int    `json:"userID"`
			Address string `json:"address"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			errors.Log(err)
			Fail(w, r, http.StatusBadRequest, "invalid payload")
			return
		}

		email, err := mail.ParseAddress(payload.Address)
		if err != nil {
			Fail(w, r, http.StatusBadRequest, "invalid email address")
			return
		}

		if payload.UserID < 1 {
			Fail(w, r, http.StatusBadRequest, "invalid user ID")
			return
		}

		emailID, err := es.Add(r.Context(), payload.UserID, email.Address)
		if err != nil {
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, struct {
			EmailID int `json:"emailID"`
		}{emailID})
	}
}

func DeleteEmailHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailID, err := strconv.ParseUint(chi.URLParam(r, "emailID"), 10, 64)
		if err != nil || emailID <= 0 {
			Fail(w, r, http.StatusBadRequest, "invalid email ID")
			return
		}

		err = es.Delete(r.Context(), int(emailID))
		if err != nil {
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, nil)
	}
}

func ListEmailsHandle(es iface.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()["userID"]
		if len(params) == 0 {
			Fail(w, r, http.StatusBadRequest, "missing URL query userID")
			return
		}

		userID, err := strconv.ParseUint(params[0], 10, 64)
		if err != nil || userID == 0 {
			Fail(w, r, http.StatusBadRequest, "invalid URL query userID")
			return
		}

		emails, err := es.ByUserID(r.Context(), int(userID))
		if err != nil {
			errors.Log(err)
			Fail(w, r, http.StatusInternalServerError, "service failed")
			return
		}

		JSON(w, r, map[string]interface{}{
			"emails": emails,
		})
	}
}
