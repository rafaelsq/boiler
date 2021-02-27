package mutation

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"boiler/cmd/server/internal/graphql/entity"
	lentity "boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/service/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockInterface(ctrl)

	m := NewMutation(service)

	ctx := context.TODO()

	// succeed
	{
		name := "name"
		password := "pass"

		service.EXPECT().AddUser(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, u *lentity.User) error {
				u.ID = 1
				return nil
			})

		r, err := m.AddUser(ctx, entity.AddUserInput{
			Name:     name,
			Password: password,
		})

		assert.Nil(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "1", r.User.ID)
	}

	// fails if service fails
	{
		name := "name"
		password := "pass"

		service.EXPECT().AddUser(ctx, gomock.Any()).Return(fmt.Errorf("opz"))

		r, err := m.AddUser(ctx, entity.AddUserInput{
			Name:     name,
			Password: password,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, r)
	}
}

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockInterface(ctrl)

	m := NewMutation(service)

	ctx := context.TODO()

	// succeed
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, e *lentity.Email) error {
				e.ID = 1
				return nil
			})

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Nil(t, err)
		assert.Equal(t, "1", u.Email.ID)
	}

	// fails if userID is invalid
	{
		address := "email@email.com"
		userID := "0"

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid userID")
		assert.Nil(t, u)
	}

	// fails if email is invalid
	{
		address := "email"
		userID := "1"

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  userID,
			Address: address,
		})
		assert.Equal(t, err.Error(), "input: invalid email address")
		assert.Nil(t, u)
	}

	// fails if service fails with duplicated
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, gomock.Any()).Return(errors.ErrAlreadyExists)

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Equal(t, err.Error(), fmt.Sprintf("input: %v", errors.ErrAlreadyExists))
		assert.Nil(t, u)
	}

	// fails if service fails
	{
		address := "email@email.com"
		userID := int64(12)

		service.EXPECT().AddEmail(ctx, gomock.Any()).Return(fmt.Errorf("opz"))

		u, err := m.AddEmail(ctx, entity.AddEmailInput{
			UserID:  strconv.FormatInt(userID, 10),
			Address: address,
		})
		assert.Equal(t, err.Error(), "service failed")
		assert.Nil(t, u)
	}
}
