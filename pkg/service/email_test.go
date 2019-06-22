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

func TestEmailAddService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockEmailRepository(ctrl)

	srv := service.NewEmail(m)

	ID := 13
	userID := 99
	address := "contact@example.com"

	ctx := context.Background()
	m.
		EXPECT().
		Add(ctx, userID, address).
		Return(ID, nil)

	ID, err := srv.Add(ctx, userID, address)
	assert.Nil(t, err)
	assert.Equal(t, ID, ID)
}

func TestEmailDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockEmailRepository(ctrl)

	srv := service.NewEmail(m)

	ID := 13

	ctx := context.Background()
	m.
		EXPECT().
		Delete(ctx, ID).
		Return(nil)

	err := srv.Delete(ctx, ID)
	assert.Nil(t, err)
}

func TestEmailByUserIDService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockEmailRepository(ctrl)

	srv := service.NewEmail(m)

	ID := 13
	userID := 99
	address := "contact@example.com"

	ctx := context.Background()
	m.
		EXPECT().
		ByUserID(ctx, userID).
		Return([]*entity.Email{{
			ID:      ID,
			UserID:  userID,
			Address: address,
		}}, nil)

	es, err := srv.ByUserID(ctx, userID)
	assert.Nil(t, err)
	assert.Len(t, es, 1)
	assert.Equal(t, es[0].ID, ID)
}
