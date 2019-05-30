package email_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/rafaelsq/boiler/pkg/repository/email"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	userID := 12
	address := "contact@example.com"

	r := email.New(m)
	m.EXPECT().AddEmail(gomock.Any(), userID, address).Return(1, nil)

	ID, err := r.Add(context.Background(), userID, address)
	assert.Nil(t, err)
	assert.Equal(t, ID, 1)
}

func TestByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	userID := 12
	emails := []*entity.Email{{ID: 1, Address: "contact@example.com"}}

	r := email.New(m)
	m.EXPECT().EmailsByUserID(gomock.Any(), userID).Return(emails, nil)

	es, err := r.ByUserID(context.Background(), userID)
	assert.Nil(t, err)
	assert.Len(t, es, len(emails))
	for i, e := range es {
		assert.Equal(t, e.ID, emails[i].ID)
	}
}
