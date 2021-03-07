package rest

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/service"
	"boiler/pkg/store"

	"github.com/go-chi/chi"
)

func New(srv service.Interface, resp Resp) Handle {
	return Handle{srv, resp}
}

type Handle struct {
	service service.Interface
	resp    Resp
}

// AddUser handle an AddUser request
func (h *Handle) AddUser(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.resp.Fail(w, r, errors.AddCode(errors.ErrBadRequest, "invalid_payload"))
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if len(payload.Name) == 0 {
		h.resp.Fail(w, r, errors.ErrInvalidName)
		return
	}

	user := entity.User{
		Name:     payload.Name,
		Password: payload.Password,
	}

	err = h.service.AddUser(r.Context(), &user)
	if err != nil {
		h.resp.Failf(w, r, "could not add user; %w", err)
		return
	}

	h.resp.JSON(w, r, map[string]interface{}{
		"user_id": user.ID,
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
			h.resp.Fail(w, r, errors.ErrInvalidLimit)
			return
		}
	}

	users := make([]entity.User, 0)
	err := h.service.FilterUsers(r.Context(), store.FilterUsers{Limit: uint(limit)}, &users)
	if err != nil {
		h.resp.Failf(w, r, "could not filter users; %w", err)
		return
	}

	h.resp.JSON(w, r, map[string]interface{}{
		"users": users,
	})
}

// DeleteUser handle an DeleteUser request
func (h *Handle) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID == 0 {
		h.resp.Fail(w, r, errors.ErrInvalidUserID)
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		h.resp.Failf(w, r, "could not delete user; %w", err)
		return
	}

	h.resp.JSON(w, r, nil)
}

// GetUser handle an GetUser request
func (h *Handle) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID == 0 {
		h.resp.Fail(w, r, errors.ErrInvalidUserID)
		return
	}

	user := new(entity.User)
	err = h.service.GetUserByID(r.Context(), userID, user)
	if err != nil {
		h.resp.Failf(w, r, "could not get user; %w", err)
		return
	}

	h.resp.JSON(w, r, map[string]interface{}{
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
		h.resp.Fail(w, r, errors.AddCode(errors.ErrBadRequest, "invalid_payload"))
		return
	}

	emailAddress, err := mail.ParseAddress(payload.Address)
	if err != nil {
		h.resp.Fail(w, r, errors.ErrInvalidEmailAddress)
		return
	}

	if payload.UserID < 1 {
		h.resp.Fail(w, r, errors.ErrInvalidID)
		return
	}

	email := entity.Email{
		UserID:  payload.UserID,
		Address: emailAddress.Address,
	}

	err = h.service.AddEmail(r.Context(), &email)
	if err != nil {
		h.resp.Failf(w, r, "could not add email; %w", err)
		return
	}

	h.resp.JSON(w, r, struct {
		EmailID int64 `json:"email_id"`
	}{email.ID})
}

// DeleteEmail handle an DeleteEmail request
func (h *Handle) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	emailID, err := strconv.ParseInt(chi.URLParam(r, "emailID"), 10, 64)
	if err != nil || emailID <= 0 {
		h.resp.Fail(w, r, errors.ErrInvalidEmailID)
		return
	}

	err = h.service.DeleteEmail(r.Context(), emailID)
	if err != nil {
		h.resp.Failf(w, r, "could not delete email; %w", err)
		return
	}

	h.resp.JSON(w, r, nil)
}

// ListEmails handle an ListEmails request
func (h *Handle) ListEmails(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()["user_id"]
	if len(params) == 0 {
		h.resp.Fail(w, r, errors.AddCode(errors.ErrBadRequest, "missing_query_user_id"))
		return
	}

	userID, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil || userID == 0 {
		h.resp.Fail(w, r, errors.ErrInvalidUserID)
		return
	}

	emails := make([]entity.Email, 0)
	err = h.service.FilterEmails(r.Context(), store.FilterEmails{UserID: userID}, &emails)
	if err != nil {
		h.resp.Failf(w, r, "could not filter email; %w", err)
		return
	}

	h.resp.JSON(w, r, map[string]interface{}{
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
		h.resp.Fail(w, r, errors.AddCode(errors.ErrBadRequest, "invalid_payload"))
		return
	}

	var user *entity.User
	var token *string

	err = h.service.AuthUser(r.Context(), payload.Email, payload.Password, user, token)
	if err != nil {
		h.resp.Failf(w, r, "could not auth user; %w", err)
		return
	}

	h.resp.JSON(w, r, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
