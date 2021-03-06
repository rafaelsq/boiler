package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"boiler/cmd/server/internal/rest"
	"boiler/cmd/server/internal/router"
	"boiler/pkg/entity"
	"boiler/pkg/service/mock"
	"boiler/pkg/store"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddUserHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		user := &entity.User{Name: "John"}
		m.EXPECT().AddUser(gomock.Any(), user).DoAndReturn(func(_ context.Context, u *entity.User) error {
			u.ID = 4
			return nil
		})

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/users", h.AddUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString("{\"name\":\"John\"}")
		res, err := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 0)
		assert.Nil(t, err)
	}

	// fail if payload is invalid
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/users", h.AddUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Post(
			fmt.Sprintf("%s/users?debug", ts.URL),
			"application/json",
			bytes.NewBufferString("{\"invalid}"),
		)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_payload", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fail if name is empty
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/users", h.AddUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Post(
			fmt.Sprintf("%s/users?debug", ts.URL),
			"application/json",
			bytes.NewBufferString("{\"name\":\"\"}"),
		)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_name", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		myErr := fmt.Errorf("opz")

		user := &entity.User{Name: "John"}
		m.EXPECT().AddUser(gomock.Any(), user).Return(myErr)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/users", h.AddUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString("{\"name\":\"John\"}")
		res, err := http.Post(fmt.Sprintf("%s/users?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 1)
		assert.Equal(t, "internal_server_error", resp.Error.Codes[0])
		assert.Equal(t, "could not add user; opz", resp.Error.Msg)
		assert.Nil(t, err)
	}
}

func TestDeleteUserHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		user := &entity.User{ID: 4, Name: "John"}
		m.EXPECT().DeleteUser(gomock.Any(), user.ID).Return(nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/users/{userID:[0-9]+}", h.DeleteUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d?debug", ts.URL, user.ID), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		res.Body.Close()
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)
	}

	// fails if invalid userID
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/users/{userID:[0-9]+}", h.DeleteUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/0?debug", ts.URL), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 3)
		assert.Equal(t, "invalid_user_id", resp.Error.Codes[0])
		assert.Equal(t, "invalid_id", resp.Error.Codes[1])
		assert.Equal(t, "bad_request", resp.Error.Codes[2])
		assert.Nil(t, err)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		user := &entity.User{ID: 4, Name: "John"}
		m.EXPECT().DeleteUser(gomock.Any(), user.ID).Return(fmt.Errorf("opz"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/users/{userID:[0-9]+}", h.DeleteUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d?debug", ts.URL, user.ID), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 1)
		assert.Equal(t, "internal_server_error", resp.Error.Codes[0])
		assert.Equal(t, "could not delete user; opz", resp.Error.Msg)
	}
}

func TestUsersHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		user := &entity.User{ID: 4, Name: "John Doe"}
		users := []entity.User{}
		m.EXPECT().
			FilterUsers(gomock.Any(), store.FilterUsers{Limit: 3}, &users).
			DoAndReturn(func(_ context.Context, _ store.FilterUsers, us *[]entity.User) error {
				*us = append(*us, *user)
				return nil
			})

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/users", h.ListUsers)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users?limit=3", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp struct{ Users []*entity.User }
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Len(t, resp.Users, 1)
		assert.Equal(t, resp.Users[0].ID, user.ID)
		assert.Equal(t, resp.Users[0].Name, user.Name)
	}

	// fail if invalid limit
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/users", h.ListUsers)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users?limit=a", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_limit", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		users := []entity.User{}
		m.EXPECT().FilterUsers(gomock.Any(), store.FilterUsers{Limit: 100}, &users).Return(fmt.Errorf("not working"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/users", h.ListUsers)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users?debug", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 1)
		assert.Equal(t, "internal_server_error", resp.Error.Codes[0])
		assert.Equal(t, "could not filter users; not working", resp.Error.Msg)
	}
}

func TestUserHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		m.EXPECT().
			GetUserByID(gomock.Any(), int64(4), gomock.Any()).
			DoAndReturn(func(_ context.Context, id int64, u *entity.User) error {
				u.ID = 4
				u.Name = "John Doe"
				return nil
			})

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/user/{userID:[0-9]+}", h.GetUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/4", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var u struct{ User *entity.User }
		err = json.NewDecoder(res.Body).Decode(&u)
		res.Body.Close()
		assert.Nil(t, err)
		assert.NotNil(t, u)

		assert.Equal(t, u.User.ID, int64(4))
		assert.Equal(t, u.User.Name, "John Doe")
	}

	// fail - bad-request
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/user/{userID:[0-9]+}", h.GetUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/0", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		m.EXPECT().
			GetUserByID(gomock.Any(), int64(4), gomock.Any()).
			Return(fmt.Errorf("failed"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/user/{userID:[0-9]+}", h.GetUser)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/4", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		res.Body.Close()
	}
}

func TestAddEmailHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		email := entity.Email{
			UserID:  12,
			Address: "example@email.com",
		}

		m.EXPECT().
			AddEmail(gomock.Any(), &email).
			DoAndReturn(func(_ context.Context, e *entity.Email) error {
				e.ID = 5
				return nil
			})

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/emails", h.AddEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString(fmt.Sprintf("{\"user_id\":%d,\"address\":\"%s\"}", email.UserID, email.Address))

		res, err := http.Post(fmt.Sprintf("%s/emails?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		j := struct {
			EmailID int `json:"email_id"`
		}{}
		err = json.NewDecoder(res.Body).Decode(&j)
		res.Body.Close()
		assert.Nil(t, err)
		assert.Equal(t, 5, j.EmailID)
	}

	// fail with invalid payload
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/emails", h.AddEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString("{invalid-payload}")

		res, err := http.Post(fmt.Sprintf("%s/emails?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_payload", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fail with invalid email address
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/emails", h.AddEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString("{\"userID\":12,\"address\":\"invalid-email\"}")

		res, err := http.Post(fmt.Sprintf("%s/emails?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_email_address", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fail with invalid user ID
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/emails", h.AddEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString("{\"userID\":0,\"address\":\"example@email.com\"}")

		res, err := http.Post(fmt.Sprintf("%s/emails?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 2)
		assert.Equal(t, "invalid_id", resp.Error.Codes[0])
		assert.Equal(t, "bad_request", resp.Error.Codes[1])
		assert.Nil(t, err)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		email := entity.Email{
			UserID:  12,
			Address: "example@email.com",
		}
		myErr := fmt.Errorf("fails")

		m.EXPECT().
			AddEmail(gomock.Any(), &email).
			Return(myErr)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Post("/emails", h.AddEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBufferString(fmt.Sprintf("{\"user_id\":%d,\"address\":\"%s\"}", email.UserID, email.Address))

		res, err := http.Post(fmt.Sprintf("%s/emails?debug", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 1)
		assert.Equal(t, "internal_server_error", resp.Error.Codes[0])
		assert.Equal(t, "could not add email; fails", resp.Error.Msg)
		assert.Nil(t, err)
	}
}

func TestDeleteEmailHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		emailID := int64(12)

		m.EXPECT().
			DeleteEmail(gomock.Any(), emailID).
			Return(nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/emails/{emailID:[0-9]+}", h.DeleteEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/emails/%d?debug", ts.URL, emailID), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)

		res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK)
	}

	// fails if emailID is invalid
	{
		m := mock.NewMockInterface(ctrl)

		emailID := int64(0)

		// m.EXPECT().
		// 	Delete(gomock.Any(), emailID).
		// 	Return(nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/emails/{emailID:[0-9]+}", h.DeleteEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/emails/%d?debug", ts.URL, emailID), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		emailID := int64(1)

		m.EXPECT().
			DeleteEmail(gomock.Any(), emailID).
			Return(fmt.Errorf("opz"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Delete("/emails/{emailID:[0-9]+}", h.DeleteEmail)

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/emails/%d?debug", ts.URL, emailID), nil)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		var resp rest.ErrResponse
		assert.Nil(t, json.NewDecoder(res.Body).Decode(&resp))
		res.Body.Close()

		assert.Len(t, resp.Error.Codes, 1)
		assert.Equal(t, "internal_server_error", resp.Error.Codes[0])
		assert.Equal(t, "could not delete email; opz", resp.Error.Msg)
		assert.Nil(t, err)
	}
}

func TestEmailsHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)

		user := entity.User{ID: 4, Name: "John Doe"}
		emails := []entity.Email{
			{ID: 2, Address: "contact@example.com"},
			{ID: 3, Address: "devs@example.com"},
		}

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{UserID: user.ID}, gomock.Any()).
			DoAndReturn(func(_ context.Context, f store.FilterEmails, es *[]entity.Email) error {
				*es = emails
				return nil
			})

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/emails", h.ListEmails)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails?debug&user_id=%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var rEmails struct{ Emails []*entity.Email }
		err = json.NewDecoder(res.Body).Decode(&rEmails)
		res.Body.Close()
		assert.Nil(t, err)
		assert.NotNil(t, rEmails)

		assert.Len(t, rEmails.Emails, len(emails))
		for i, e := range rEmails.Emails {
			assert.Equal(t, e.ID, emails[i].ID)
			assert.Equal(t, e.Address, emails[i].Address)
		}
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)

		user := entity.User{ID: 4, Name: "John Doe"}

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{UserID: user.ID}, gomock.Any()).
			Return(fmt.Errorf("failed"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/emails", h.ListEmails)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails?user_id=%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		res.Body.Close()
	}

	// fail - bad request
	{
		m := mock.NewMockInterface(ctrl)

		h := rest.New(m, new(rest.DefaultResp))
		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)
		r.Get("/emails", h.ListEmails)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails?user_id=0", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}

	// fails if not userID is provided
	{
		m := mock.NewMockInterface(ctrl)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r, nil, m)

		h := rest.New(m, new(rest.DefaultResp))
		r.Get("/emails", h.ListEmails)

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}
}
