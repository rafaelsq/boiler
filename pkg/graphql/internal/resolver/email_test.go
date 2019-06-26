package resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	gentity "github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/resolver"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestEmailID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	r := resolver.NewEmail(m)

	emailID, err := r.ID(context.TODO(), &gentity.Email{ID: 4})
	assert.Nil(t, err)
	assert.Equal(t, 4, emailID)
}

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

		u, err := r.User(context.TODO(), &gentity.Email{Address: email.Address})
		assert.Nil(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, u.ID, user.ID)
	}

	// fails if service fails
	{
		email := &entity.Email{ID: 4, Address: "a@b.c"}

		m := mock.NewMockService(ctrl)
		r := resolver.NewEmail(m)

		m.EXPECT().
			GetUserByEmail(gomock.Any(), email.Address).
			Return(nil, fmt.Errorf("opz"))

		u, err := r.User(context.TODO(), &gentity.Email{Address: email.Address})
		assert.Nil(t, u)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}
}
