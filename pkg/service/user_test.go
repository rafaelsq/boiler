package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestUserAddService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockUserRepository(ctrl)

	srv := service.NewUser(m)

	userID := 99
	name := "name"

	ctx := context.Background()
	m.
		EXPECT().
		Add(ctx, name).
		Return(userID, nil)

	id, err := srv.Add(ctx, name)
	assert.Nil(t, err)
	assert.Equal(t, userID, id)
}

func TestUserDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockUserRepository(ctrl)

	srv := service.NewUser(m)

	userID := 99

	ctx := context.Background()
	m.
		EXPECT().
		Delete(ctx, userID).
		Return(nil)

	err := srv.Delete(ctx, userID)
	assert.Nil(t, err)
}

func TestUserListService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockUserRepository(ctrl)

	srv := service.NewUser(m)

	userID := 99
	name := "name"

	ctx := context.Background()
	m.
		EXPECT().
		List(ctx, uint(10)).
		Return([]*entity.User{{
			ID:   userID,
			Name: name,
		}}, nil)

	vs, err := srv.List(ctx, 10)
	assert.Nil(t, err)
	assert.Len(t, vs, 1)
	assert.Equal(t, vs[0].ID, userID)
	assert.Equal(t, vs[0].Name, name)
}

func TestUserByIDService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockUserRepository(ctrl)

	srv := service.NewUser(m)

	userID := 99
	name := "userName"

	ctx := context.Background()
	m.
		EXPECT().
		ByID(ctx, gomock.Eq(userID)).
		Return(&entity.User{
			ID:   userID,
			Name: name,
		}, nil)

	v, err := srv.ByID(ctx, userID)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, v.Name, name)
}

func TestUserByEmailService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockUserRepository(ctrl)

	srv := service.NewUser(m)

	userID := 99
	name := "userName"
	email := "contact@example.com"

	ctx := context.Background()
	m.
		EXPECT().
		ByEmail(ctx, gomock.Eq(email)).
		Return(&entity.User{
			ID:   userID,
			Name: name,
		}, nil)

	v, err := srv.ByEmail(ctx, email)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, v.Name, name)
}
