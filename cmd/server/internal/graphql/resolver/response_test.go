package resolver_test

import (
	"context"
	"testing"

	"boiler/cmd/server/internal/graphql/entity"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/errors"
	"boiler/pkg/service/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResponseUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)
	r := resolver.NewResponse(m)
	_, err := r.User(context.TODO(), &entity.UserResponse{
		User: &entity.User{ID: ""},
	})
	assert.Equal(t, err, errors.ErrInvalidID)
}

func TestResponseEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)
	r := resolver.NewResponse(m)
	_, err := r.Email(context.TODO(), &entity.EmailResponse{
		Email: &entity.Email{ID: ""},
	})
	assert.Equal(t, err, errors.ErrInvalidID)
}

func TestAuthUserResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)
	r := resolver.NewAuthUserResponse(m)
	_, err := r.User(context.TODO(), &entity.AuthUserResponse{
		User: &entity.User{ID: ""},
	})
	assert.Equal(t, err, errors.ErrInvalidID)
}
