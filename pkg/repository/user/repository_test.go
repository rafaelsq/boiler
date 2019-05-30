package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	r := user.New(m)
	m.EXPECT().Users(gomock.Any()).Return([]*entity.User{
		{ID: 1}, {ID: 2},
	}, nil)

	users, err := r.List(context.Background())
	assert.Nil(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, users[0].ID, 1)
	assert.Equal(t, users[1].ID, 2)
}

func TestByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	userID := 3

	r := user.New(m)
	m.EXPECT().UserByID(gomock.Any(), userID).Return(&entity.User{
		ID: userID,
	}, nil)

	u, err := r.ByID(context.Background(), userID)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, userID)
}

func TestByMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	userID := 3
	email := "contact@example.com"

	r := user.New(m)
	m.EXPECT().UserByEmail(gomock.Any(), email).Return(&entity.User{
		ID: userID,
	}, nil)

	u, err := r.ByEmail(context.Background(), email)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, userID)
}
