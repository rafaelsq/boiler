package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/rafaelsq/boiler/pkg/rest"
	"github.com/rafaelsq/boiler/pkg/router"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestAddUserHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		user := &entity.User{ID: 4, Name: "John"}
		m.EXPECT().Add(gomock.Any(), user.Name).Return(user.ID, nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Post("/users/add", rest.AddUserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBuffer([]byte("{\"name\":\"John\"}"))
		res, err := http.Post(fmt.Sprintf("%s/users/add", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var rm map[string]int
		_ = json.NewDecoder(res.Body).Decode(&rm)
		res.Body.Close()
		assert.Equal(t, rm["UserID"], user.ID)
		assert.Nil(t, err)
	}

	// fail if payload is invalid
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Post("/users/add", rest.AddUserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Post(
			fmt.Sprintf("%s/users/add", ts.URL),
			"application/json",
			bytes.NewBuffer([]byte("{\"invalid}")),
		)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(b), "could not parse payload")
		res.Body.Close()
	}

	// fail if name is empty
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Post("/users/add", rest.AddUserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Post(
			fmt.Sprintf("%s/users/add", ts.URL),
			"application/json",
			bytes.NewBuffer([]byte("{\"name\":\"\"}")),
		)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(b), "empty name")
		res.Body.Close()
	}

	// fail if service fail
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)
		myErr := fmt.Errorf("opz")

		user := &entity.User{ID: 4, Name: "John"}
		m.EXPECT().Add(gomock.Any(), user.Name).Return(0, myErr)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Post("/users/add", rest.AddUserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		body := bytes.NewBuffer([]byte("{\"name\":\"John\"}"))
		res, err := http.Post(fmt.Sprintf("%s/users/add", ts.URL), "application/json", body)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)

		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(b), "service fail")
		res.Body.Close()
	}
}

func TestUsersHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		user := &entity.User{ID: 4, Name: "John Doe"}
		m.EXPECT().List(gomock.Any(), uint(3)).Return([]*entity.User{user}, nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/users", rest.UsersHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users?limit=3", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var users []*entity.User
		err = json.NewDecoder(res.Body).Decode(&users)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Len(t, users, 1)
		assert.Equal(t, users[0].ID, user.ID)
		assert.Equal(t, users[0].Name, user.Name)
	}

	// fail if invalid limit
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/users", rest.UsersHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users?limit=a", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)

		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(b), "invalid limit \"a\"")
		res.Body.Close()
	}

	// fail if service fail
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		user := &entity.User{ID: 4, Name: "John Doe"}
		m.EXPECT().List(gomock.Any(), uint(100)).Return([]*entity.User{user}, fmt.Errorf("not working"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/users", rest.UsersHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/users", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)

		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(b), "service fail")
		res.Body.Close()
	}
}

func TestUserHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		user := &entity.User{ID: 4, Name: "John Doe"}

		m.EXPECT().
			ByID(gomock.Any(), user.ID).
			Return(user, nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/user/{userID:[0-9]+}", rest.UserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var u *entity.User
		err = json.NewDecoder(res.Body).Decode(&u)
		res.Body.Close()
		assert.Nil(t, err)
		assert.NotNil(t, u)

		assert.Equal(t, u.ID, user.ID)
		assert.Equal(t, u.Name, user.Name)
	}

	// fail - bad-request
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/user/{userID:[0-9]+}", rest.UserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/0", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}

	// fail - service failed
	{
		m := mock.NewMockUserRepository(ctrl)
		us := service.NewUser(m)

		user := &entity.User{ID: 4, Name: "John Doe"}

		m.EXPECT().
			ByID(gomock.Any(), user.ID).
			Return(nil, fmt.Errorf("failed"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/user/{userID:[0-9]+}", rest.UserHandle(us))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/user/%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		res.Body.Close()
	}
}

func TestEmailsHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockEmailRepository(ctrl)
		es := service.NewEmail(m)

		user := entity.User{ID: 4, Name: "John Doe"}
		emails := []*entity.Email{
			{ID: 2, Address: "contact@example.com"},
			{ID: 3, Address: "devs@example.com"},
		}

		m.EXPECT().
			ByUserID(gomock.Any(), user.ID).
			Return(emails, nil)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/emails/{userID:[0-9]+}", rest.EmailsHandle(es))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails/%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var rEmails []*entity.Email
		err = json.NewDecoder(res.Body).Decode(&rEmails)
		res.Body.Close()
		assert.Nil(t, err)
		assert.NotNil(t, rEmails)

		assert.Len(t, rEmails, len(emails))
		for i, e := range rEmails {
			assert.Equal(t, e.ID, emails[i].ID)
			assert.Equal(t, e.Address, emails[i].Address)
		}
	}

	// fail - service failed
	{
		m := mock.NewMockEmailRepository(ctrl)
		es := service.NewEmail(m)

		user := entity.User{ID: 4, Name: "John Doe"}

		m.EXPECT().
			ByUserID(gomock.Any(), user.ID).
			Return(nil, fmt.Errorf("failed"))

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/emails/{userID:[0-9]+}", rest.EmailsHandle(es))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails/%d", ts.URL, user.ID))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		res.Body.Close()
	}

	// fail - bad request
	{
		m := mock.NewMockEmailRepository(ctrl)
		es := service.NewEmail(m)

		r := chi.NewRouter()
		router.ApplyMiddlewares(r)
		r.Get("/emails/{userID:[0-9]+}", rest.EmailsHandle(es))

		ts := httptest.NewServer(r)
		defer ts.Close()

		res, err := http.Get(fmt.Sprintf("%s/emails/0", ts.URL))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		res.Body.Close()
	}
}
