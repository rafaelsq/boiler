package mutation

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mock.NewMockUserService(ctrl)

	m := NewMutation(us, nil)

	ctx := context.TODO()

	// succeed
	{
		name := "name"

		us.EXPECT().Add(ctx, name).Return(1, nil)

		u, err := m.AddUser(ctx, entity.AddUserInput{
			Name: name,
		})
		assert.Nil(t, err)
		assert.NotNil(t, u)
	}

	// fails if service fails
	{
		name := "name"

		us.EXPECT().Add(ctx, name).Return(0, fmt.Errorf("opz"))

		u, err := m.AddUser(ctx, entity.AddUserInput{
			Name: name,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, u)
	}
}

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := mock.NewMockEmailService(ctrl)

	m := NewMutation(nil, es)

	ctx := context.TODO()

	// succeed
	{
		address := "email@email.com"
		userID := 12

		es.EXPECT().Add(ctx, userID, address).Return(1, nil)

		u, err := m.AddMail(ctx, entity.AddMailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Nil(t, err)
		assert.Equal(t, u.ID, userID)
	}

	// fails if userID is invalid
	{
		address := "email@email.com"
		userID := 0

		u, err := m.AddMail(ctx, entity.AddMailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid userID")
		assert.Nil(t, u)
	}

	// fails if email is invalid
	{
		address := "email"
		userID := 1

		u, err := m.AddMail(ctx, entity.AddMailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid email address")
		assert.Nil(t, u)
	}

	// fails if service fails
	{
		address := "email@email.com"
		userID := 12

		es.EXPECT().Add(ctx, userID, address).Return(0, fmt.Errorf("opz"))

		u, err := m.AddMail(ctx, entity.AddMailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, u)
	}
}
