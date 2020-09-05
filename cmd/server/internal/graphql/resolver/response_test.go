package resolver_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"boiler/cmd/server/internal/graphql/entity"
	"boiler/cmd/server/internal/graphql/resolver"
	"boiler/pkg/iface"
	"boiler/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestResponseUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	r := resolver.NewResponse(m)
	_, err := r.User(context.TODO(), &entity.UserResponse{
		User: &entity.User{ID: ""},
	})
	assert.Equal(t, err, iface.ErrInvalidID)
}

func TestResponseEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	r := resolver.NewResponse(m)
	_, err := r.Email(context.TODO(), &entity.EmailResponse{
		Email: &entity.Email{ID: ""},
	})
	assert.Equal(t, err, iface.ErrInvalidID)
}
