package resolver_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	gentity "github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestEmailEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		email := &entity.Email{Address: "a@b.c"}
		user := &entity.User{ID: 4, Name: "John Doe"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), email.Address).
			Return(user, nil)

		u, err := r.User(ctxDebug, &gentity.Email{Address: email.Address})
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, u.ID, strconv.Itoa(user.ID))
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 4, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), email.Address).
			Return(nil, fmt.Errorf("opz"))

		u, err := r.User(ctxDebug, &gentity.Email{Address: email.Address})
		assert.Nil(t, u)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestEmailEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		email := &entity.Email{ID: 5, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), iface.FilterEmails{
				EmailID: email.ID,
			}).
			Return([]*entity.Email{email}, nil)

		e, err := r.Email(ctxDebug, strconv.Itoa(email.ID))
		assert.Nil(t, err)
		assert.NotNil(t, e)
		assert.Equal(t, e.ID, strconv.Itoa(email.ID))
	}

	// fails if invalid ID
	{
		email := &entity.Email{ID: 0, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		e, err := r.Email(ctxDebug, strconv.Itoa(email.ID))
		assert.Nil(t, e)
		assert.Equal(t, iface.ErrInvalidID, err)
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 5, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), iface.FilterEmails{
				EmailID: email.ID,
			}).
			Return(nil, errors.New("err"))

		e, err := r.Email(ctxDebug, strconv.Itoa(email.ID))
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "err")
		assert.Nil(t, e)
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 5, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), iface.FilterEmails{
				EmailID: email.ID,
			}).
			Return([]*entity.Email{}, nil)

		e, err := r.Email(ctxDebug, strconv.Itoa(email.ID))
		assert.Equal(t, err, iface.ErrNotFound)
		assert.Nil(t, e)
	}
}
