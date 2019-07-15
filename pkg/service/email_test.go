package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/mock"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/stretchr/testify/assert"
)

type tx struct{}

func (*tx) Commit() error { return nil }

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	ID := 13
	userID := 99
	address := "contact@example.com"

	ctx := context.Background()

	// succeed
	{
		db, mdb, err := sqlmock.New()
		defer db.Close()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		m.EXPECT().Tx().Return(tx, err)
		m.
			EXPECT().
			AddEmail(ctx, gomock.Any(), userID, address).
			Return(ID, nil)

		mdb.ExpectCommit()

		ID, err = srv.AddEmail(ctx, userID, address)
		assert.Nil(t, err)
		assert.Equal(t, ID, ID)
	}

	// fails if Tx fails
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("opz"))

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not begin transaction; opz")
		assert.Equal(t, 0, id)
	}

	// fails if service fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, userID, address).
			Return(0, fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollback")
		assert.Equal(t, 0, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if service fails and rollback fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, userID, address).
			Return(0, fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollbackerr; rollback")
		assert.Equal(t, 0, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if commit fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, userID, address).
			Return(ID, nil)

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; commit failed")
		assert.Equal(t, 0, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	ID := 13

	ctx := context.Background()
	m.
		EXPECT().
		DeleteEmail(ctx, ID).
		Return(nil)

	err := srv.DeleteEmail(ctx, ID)
	assert.Nil(t, err)
}

func TestFilterEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	ID := 13
	userID := 99
	address := "contact@example.com"
	filter := iface.FilterEmails{UserID: userID}
	ctx := context.Background()
	m.
		EXPECT().
		FilterEmails(ctx, filter).
		Return([]*entity.Email{{
			ID:      ID,
			UserID:  userID,
			Address: address,
		}}, nil)

	es, err := srv.FilterEmails(ctx, filter)
	assert.Nil(t, err)
	assert.Len(t, es, 1)
	assert.Equal(t, es[0].ID, ID)
}
