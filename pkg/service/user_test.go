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

func TestUserService(t *testing.T) {
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
