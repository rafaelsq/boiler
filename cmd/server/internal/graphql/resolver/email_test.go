package resolver_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	gentity "boiler/cmd/server/internal/graphql/entity"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/entity"
	"boiler/pkg/errors"
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
		emailAddress := "a@b.c"

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), emailAddress, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, u *entity.User) error {
				u.ID = 4
				u.Name = "John Doe"
				return nil
			})

		u, err := r.User(ctxDebug, &gentity.Email{Address: emailAddress})
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, u.ID, "4")
		assert.Equal(t, u.Name, "John Doe")
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 4, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), email.Address, gomock.Any()).
			Return(fmt.Errorf("opz"))

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
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{
				EmailID: 5,
			}, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filter store.FilterEmails, emails *[]entity.Email) error {
				*emails = append(*emails, entity.Email{ID: 5, Address: "a@b.c"})
				return nil
			})

		e, err := r.Email(ctxDebug, "5")
		assert.Nil(t, err)
		assert.NotNil(t, e)
		assert.Equal(t, "5", e.ID)
	}

	// fails if invalid ID
	{
		email := &entity.Email{ID: 0, Address: "a@b.c"}

		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		e, err := r.Email(ctxDebug, strconv.FormatInt(email.ID, 10))
		assert.Nil(t, e)
		assert.Equal(t, errors.ErrInvalidID, err)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{EmailID: 500}, gomock.Any()).
			Return(errors.New("err"))

		email, err := r.Email(ctxDebug, "500")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "err")
		assert.Nil(t, email)
	}

	// fails if service fails
	{
		m := mock.NewMockInterface(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			FilterEmails(gomock.Any(), store.FilterEmails{EmailID: 404}, gomock.Any()).
			Return(nil)

		email, err := r.Email(ctxDebug, "404")
		assert.Equal(t, err, errors.ErrNotFound)
		assert.Nil(t, email)
	}
}
