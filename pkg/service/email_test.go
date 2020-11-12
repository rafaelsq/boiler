package service_test

import (
	"context"
	"fmt"
	"testing"

	"boiler/pkg/entity"
	"boiler/pkg/iface"
	"boiler/pkg/mock"
	"boiler/pkg/service"
	"boiler/pkg/store/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var ID int64 = 13
	var userID int64 = 99
	address := "contact@example.com"

	ctx := context.Background()

	// succeed
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

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
		assert.Equal(t, 0, int(id))
	}

	// fails if service fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, userID, address).
			Return(int64(0), fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollback")
		assert.Equal(t, 0, int(id))
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if service fails and rollback fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, userID, address).
			Return(int64(0), fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		id, err := srv.AddEmail(ctx, userID, address)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollbackerr; rollback")
		assert.Equal(t, 0, int(id))
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if commit fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

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
		assert.Equal(t, 0, int(id))
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var ID int64 = 13

	ctx := context.Background()

	// succeed
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)
		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			DeleteEmail(ctx, tx, ID).
			Return(nil)

		mdb.ExpectCommit()

		err = srv.DeleteEmail(ctx, ID)
		assert.Nil(t, err)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// tx
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("tx fail"))

		err := srv.DeleteEmail(ctx, ID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not begin delete email transaction; tx fail", err.Error())
	}

	// store fail
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			DeleteEmail(ctx, tx, ID).
			Return(fmt.Errorf("opz"))
		mdb.ExpectRollback()

		err = srv.DeleteEmail(ctx, ID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not delete email; opz", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// commit fail
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)
		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			DeleteEmail(ctx, tx, ID).
			Return(nil)

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit fail"))

		err = srv.DeleteEmail(ctx, ID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not commit delete email; commit fail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// rollback fail
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			DeleteEmail(ctx, tx, ID).
			Return(fmt.Errorf("database fail"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackfail"))

		err = srv.DeleteEmail(ctx, ID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not rollback delete email; rollbackfail; database fail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestFilterEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var ID int64 = 13
	var userID int64 = 99
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
