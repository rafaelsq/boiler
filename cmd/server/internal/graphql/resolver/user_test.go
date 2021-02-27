package resolver_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	gentity "boiler/cmd/server/internal/graphql/entity"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/entity"
	"boiler/pkg/service"
	"boiler/pkg/service/mock"
	"boiler/pkg/store"
	"boiler/pkg/store/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ctxDebug context.Context

func init() {
	ctxDebug = context.WithValue(context.Background(), config.ContextKeyDebug{}, true)
}

func TestUserUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			GetUserByID(gomock.Any(), int64(4), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ int64, u *entity.User) error {
				u.ID = 4
				u.Name = "John Doe"
				return nil
			})

		u, err := r.User(ctxDebug, "4")
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, "4", u.ID)
		assert.Equal(t, "John Doe", u.Name)
	}

	// fail if invalid ID
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		u, err := r.User(ctxDebug, "fail")
		assert.Nil(t, u)
		assert.Equal(t, err, service.ErrInvalidID)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			GetUserByID(gomock.Any(), int64(4), gomock.Any()).
			Return(fmt.Errorf("opz"))

		u, err := r.User(ctxDebug, "4")
		assert.Nil(t, u)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestUserUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterUsers(gomock.Any(), store.FilterUsers{Limit: 2}, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ store.FilterUsers, users *[]entity.User) error {
				*users = append(*users, entity.User{ID: 4, Name: "John Doe"})
				return nil
			})

		users, err := r.Users(ctxDebug, 2)
		assert.Nil(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, len(users), 1)
		assert.Equal(t, "John Doe", users[0].Name)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterUsers(gomock.Any(), store.FilterUsers{Limit: 4}, gomock.Any()).
			Return(fmt.Errorf("opz"))

		users, err := r.Users(ctxDebug, 4)
		assert.Len(t, users, 0)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestUserEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{UserID: 4}, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ store.FilterEmails, emails *[]entity.Email) error {
				*emails = append(*emails, entity.Email{ID: 4, Address: "a@b.c"})
				return nil
			})

		emails, err := r.Emails(ctxDebug, &gentity.User{ID: "4"})
		assert.Nil(t, err)
		assert.NotNil(t, emails)
		assert.Equal(t, len(emails), 1)
	}

	// fail if invalid ID
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		emails, err := r.Emails(ctxDebug, &gentity.User{ID: "0"})
		assert.Nil(t, emails)
		assert.Equal(t, err, service.ErrInvalidID)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewUser(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{UserID: 2}, gomock.Any()).
			Return(fmt.Errorf("opz"))

		users, err := r.Emails(ctxDebug, &gentity.User{ID: strconv.FormatInt(2, 10)})
		assert.Nil(t, users)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}
