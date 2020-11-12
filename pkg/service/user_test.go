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

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "name"
	password := "pass"

	ctx := context.Background()

	// succeed
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
			AddUser(ctx, tx, name, gomock.Any()).
			Return(userID, nil)
		mdb.ExpectCommit()

		id, err := srv.AddUser(ctx, name, password)
		assert.Nil(t, err)
		assert.Equal(t, userID, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if Tx fails
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("opz"))

		id, err := srv.AddUser(ctx, name, password)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not begin transaction; opz")
		assert.Equal(t, 0, int(id))
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
			AddUser(ctx, tx, name, gomock.Any()).
			Return(int64(0), fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		id, err := srv.AddUser(ctx, name, password)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; rollback")
		assert.Equal(t, 0, int(id))
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
			AddUser(ctx, tx, name, gomock.Any()).
			Return(int64(0), fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		id, err := srv.AddUser(ctx, name, password)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; rollbackerr; rollback")
		assert.Equal(t, 0, int(id))
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
			AddUser(ctx, tx, name, gomock.Any()).
			Return(userID, nil)

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

		id, err := srv.AddUser(ctx, name, password)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; commit failed")
		assert.Equal(t, 0, int(id))
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99

	ctx := context.Background()

	// succeed
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
			DeleteUser(ctx, tx, userID).
			Return(nil)
		m.
			EXPECT().
			DeleteEmailsByUserID(ctx, tx, userID).
			Return(nil)
		mdb.ExpectCommit()

		err = srv.DeleteUser(ctx, userID)
		assert.Nil(t, err)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if Tx fails
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("tx fails"))

		err := srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not begin delete user transaction; tx fails")
	}

	// DeleteUser fail
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
			DeleteUser(ctx, tx, userID).
			Return(fmt.Errorf("deletefail"))
		mdb.ExpectRollback()

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not delete user; deletefail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// DeleteUser rollback fail
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
			DeleteUser(ctx, tx, userID).
			Return(fmt.Errorf("deletefail"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackfail"))

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not rollback delete user; rollbackfail; deletefail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// DeleteEmail fail
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
			DeleteUser(ctx, tx, userID).
			Return(nil)
		m.
			EXPECT().
			DeleteEmailsByUserID(ctx, tx, userID).
			Return(fmt.Errorf("deletefail"))
		mdb.ExpectRollback()

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not delete user emails; deletefail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// DeleteEmail rollback fail
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
			DeleteUser(ctx, tx, userID).
			Return(nil)
		m.
			EXPECT().
			DeleteEmailsByUserID(ctx, tx, userID).
			Return(fmt.Errorf("deletefail"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackfail"))

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not rollback delete emails by user ID; rollbackfail; deletefail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// commit fail
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
			DeleteUser(ctx, tx, userID).
			Return(nil)
		m.
			EXPECT().
			DeleteEmailsByUserID(ctx, tx, userID).
			Return(nil)
		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commitfail"))

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, "could not commit delete user; commitfail", err.Error())
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestFilterUsersID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "name"
	filter := iface.FilterUsers{Limit: 10}
	ctx := context.Background()

	// succed
	{
		m.
			EXPECT().
			FilterUsersID(ctx, filter).
			Return([]int64{userID}, nil)
		m.
			EXPECT().
			FetchUsers(ctx, userID).
			Return([]*entity.User{{
				ID:   userID,
				Name: name,
			}}, nil)

		vs, err := srv.FilterUsers(ctx, filter)
		assert.Nil(t, err)
		assert.Len(t, vs, 1)
		assert.Equal(t, vs[0].ID, userID)
		assert.Equal(t, vs[0].Name, name)
	}

	// fail
	{
		m.
			EXPECT().
			FilterUsersID(ctx, filter).
			Return(nil, fmt.Errorf("opz"))

		IDs, err := srv.FilterUsers(ctx, filter)
		assert.Nil(t, IDs)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "userName"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FetchUsers(ctx, userID).
			Return([]*entity.User{
				{
					ID:   userID,
					Name: name,
				},
			}, nil)

		v, err := srv.GetUserByID(ctx, userID)
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, v.Name, name)
	}

	// fails if storage fails
	{
		m.
			EXPECT().
			FetchUsers(ctx, userID).
			Return(nil, fmt.Errorf("opz"))

		v, err := srv.GetUserByID(ctx, userID)
		assert.Nil(t, v)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}

	// fails if no user found
	{
		m.
			EXPECT().
			FetchUsers(ctx, userID).
			Return([]*entity.User{}, nil)

		v, err := srv.GetUserByID(ctx, userID)
		assert.Nil(t, v)
		assert.Equal(t, iface.ErrNotFound, err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStore(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "userName"
	email := "contact@example.com"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FilterUsersID(ctx, iface.FilterUsers{Email: email}).
			Return([]int64{userID}, nil)
		m.
			EXPECT().
			FetchUsers(ctx, userID).
			Return([]*entity.User{
				{
					ID:   userID,
					Name: name,
				},
			}, nil)

		v, err := srv.GetUserByEmail(ctx, email)
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, v.Name, name)
	}

	// fails if storage fails
	{
		m.
			EXPECT().
			FilterUsersID(ctx, iface.FilterUsers{Email: email}).
			Return(nil, fmt.Errorf("opz"))

		v, err := srv.GetUserByEmail(ctx, email)
		assert.Nil(t, v)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}

	// fails if no user found
	{
		m.
			EXPECT().
			FilterUsersID(ctx, iface.FilterUsers{Email: email}).
			Return([]int64{}, nil)

		v, err := srv.GetUserByEmail(ctx, email)
		assert.Nil(t, v)
		assert.Equal(t, iface.ErrNotFound, err)
	}
}
