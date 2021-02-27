package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"boiler/pkg/entity"
	"boiler/pkg/service"
	"boiler/pkg/store"
	"boiler/pkg/store/config"
	"boiler/pkg/store/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

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

		email := entity.Email{
			UserID:  userID,
			Address: address,
		}
		tx, err := db.Begin()
		assert.Nil(t, err)

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, &email).
			DoAndReturn(func(_ context.Context, _ *sql.Tx, e *entity.Email) error {
				e.ID = ID
				return nil
			})

		mdb.ExpectCommit()

		err = srv.AddEmail(ctx, &email)
		assert.Nil(t, err)
		assert.Equal(t, ID, email.ID)
	}

	// fails if Tx fails
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("opz"))

		email := entity.Email{
			UserID:  userID,
			Address: address,
		}

		err := srv.AddEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not begin transaction; opz")
	}

	// fails if service fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer func() { _ = db.Close() }()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		email := entity.Email{
			UserID:  userID,
			Address: address,
		}
		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, &email).
			Return(fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		err = srv.AddEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollback")
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

		email := entity.Email{
			UserID:  userID,
			Address: address,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, &email).
			Return(fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		err = srv.AddEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; rollbackerr; rollback")
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

		email := entity.Email{
			UserID:  userID,
			Address: address,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddEmail(ctx, tx, &email).
			DoAndReturn(func(_ context.Context, _ *sql.Tx, e *entity.Email) error {
				e.ID = ID
				return nil
			})

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

		err = srv.AddEmail(ctx, &email)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add email; commit failed")
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

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
		assert.Equal(t, "could not delete email; rollbackfail; database fail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestFilterEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var ID int64 = 13
	var userID int64 = 99
	address := "contact@example.com"
	filter := store.FilterEmails{UserID: userID}
	ctx := context.Background()

	var emails []entity.Email

	m.
		EXPECT().
		FilterEmails(ctx, filter, &emails).
		DoAndReturn(func(_ context.Context, _ store.FilterEmails, es *[]entity.Email) error {
			*es = append(*es, entity.Email{
				ID:      ID,
				UserID:  userID,
				Address: address,
			})
			return nil
		})

	err := srv.FilterEmails(ctx, filter, &emails)
	assert.Nil(t, err)
	assert.Len(t, emails, 1)
	assert.Equal(t, ID, emails[0].ID)
}
