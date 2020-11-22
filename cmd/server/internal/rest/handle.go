package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"boiler/pkg/service"
	"boiler/pkg/store"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func New(srv service.Interface) Handle {
	return Handle{srv}
}

type Handle struct {
	service service.Interface
}

// AddUser handle an AddUser request
func (h *Handle) AddUser(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
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

	userID, err := h.service.AddUser(r.Context(), payload.Name, payload.Password)
	if err != nil {
		log.Error().Err(err).Msg("could not add user")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, map[string]interface{}{
		"user_id": userID,
	})
}

// ListUsers handle an ListUsers request
func (h *Handle) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.service.FilterUsers(r.Context(), store.FilterUsers{Limit: uint(limit)})
	if err != nil {
		log.Error().Err(err).Msg("could not filter users")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "service failed")
		return
	}

	JSON(w, r, map[string]interface{}{
		"users": users,
	})
}

// DeleteUser handle an DeleteUser request
func (h *Handle) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID == 0 {
		Fail(w, r, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("could not delete user")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, nil)
}

// GetUser handle an GetUser request
func (h *Handle) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID == 0 {
		Fail(w, r, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("could not get user")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, map[string]interface{}{
		"user": user,
	})
}

// AddEmail handle an AddEmail request
func (h *Handle) AddEmail(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		UserID  int64  `json:"user_id"`
		Address string `json:"address"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
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

	emailID, err := h.service.AddEmail(r.Context(), payload.UserID, email.Address)
	if err != nil {
		log.Error().Err(err).Msg("could not add email")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, struct {
		EmailID int64 `json:"email_id"`
	}{emailID})
}

// DeleteEmail handle an DeleteEmail request
func (h *Handle) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	emailID, err := strconv.ParseInt(chi.URLParam(r, "emailID"), 10, 64)
	if err != nil || emailID <= 0 {
		Fail(w, r, http.StatusBadRequest, "invalid email ID")
		return
	}

	err = h.service.DeleteEmail(r.Context(), emailID)
	if err != nil {
		log.Error().Err(err).Msg("could not delete email")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, nil)
}

// ListEmails handle an ListEmails request
func (h *Handle) ListEmails(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()["user_id"]
	if len(params) == 0 {
		Fail(w, r, http.StatusBadRequest, "missing URL query user_id")
		return
	}

	userID, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil || userID == 0 {
		Fail(w, r, http.StatusBadRequest, "invalid URL query user_id")
		return
	}

	emails, err := h.service.FilterEmails(r.Context(), store.FilterEmails{UserID: userID})
	if err != nil {
		log.Error().Err(err).Msg("could not filter email")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, map[string]interface{}{
		"emails": emails,
	})
}

// AuthUser handle an authentication request
func (h *Handle) AuthUser(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		Fail(w, r, http.StatusBadRequest, "could not parse payload")
		return
	}

	user, token, err := h.service.AuthUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		log.Error().Err(err).Msg("could not auth user")
		Fail(w, r, http.StatusInternalServerError, "service failed")
		return
	}

	JSON(w, r, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
