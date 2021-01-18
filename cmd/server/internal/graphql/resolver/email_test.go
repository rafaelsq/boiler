package resolver_test

import (
	"fmt"
	"strconv"
	"testing"

	gentity "boiler/cmd/server/internal/graphql/entity"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/service"
	"boiler/pkg/service/mock"
	"boiler/pkg/store"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// succeed
	{
		email := &entity.Email{Address: "a@b.c"}
		user := &entity.User{ID: 4, Name: "John Doe"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), email.Address).
			Return(user, nil)

		u, err := r.User(ctxDebug, &gentity.Email{Address: email.Address})
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, u.ID, strconv.FormatInt(user.ID, 10))
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 4, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
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

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{
				EmailID: email.ID,
			}).
			Return([]*entity.Email{email}, nil)

		e, err := r.Email(ctxDebug, strconv.FormatInt(email.ID, 10))
		assert.Nil(t, err)
		assert.NotNil(t, e)
		assert.Equal(t, e.ID, strconv.FormatInt(email.ID, 10))
	}

	// fails if invalid ID
	{
		email := &entity.Email{ID: 0, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		e, err := r.Email(ctxDebug, strconv.FormatInt(email.ID, 10))
		assert.Nil(t, e)
		assert.Equal(t, service.ErrInvalidID, err)
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 5, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{
				EmailID: email.ID,
			}).
			Return(nil, errors.New("err"))

		e, err := r.Email(ctxDebug, strconv.FormatInt(email.ID, 10))
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "err")
		assert.Nil(t, e)
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 5, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{
				EmailID: email.ID,
			}).
			Return([]*entity.Email{}, nil)

		e, err := r.Email(ctxDebug, strconv.FormatInt(email.ID, 10))
		assert.Equal(t, err, errors.ErrNotFound)
		assert.Nil(t, e)
	}
}
