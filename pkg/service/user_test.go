package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/service"
	"boiler/pkg/store"
	"boiler/pkg/store/config"
	"boiler/pkg/store/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

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

		user := entity.User{
			Name:     name,
			Password: password,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddUser(ctx, tx, &user).
			DoAndReturn(func(_ context.Context, _ *sql.Tx, u *entity.User) error {
				u.ID = userID
				return nil
			})
		mdb.ExpectCommit()

		err = srv.AddUser(ctx, &user)
		assert.Nil(t, err)
		assert.Equal(t, userID, user.ID)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if Tx fails
	{
		opz := fmt.Errorf("opz")
		m.EXPECT().Tx().Return(nil, opz)

		user := entity.User{
			Name:     name,
			Password: password,
		}

		err := srv.AddUser(ctx, &user)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, opz))
		assert.Equal(t, err.Error(), "could not begin transaction; opz")
	}

	// fails if service fails
	{
		db, mdb, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		mdb.ExpectBegin()

		tx, err := db.Begin()
		assert.Nil(t, err)

		user := entity.User{
			Name:     name,
			Password: password,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddUser(ctx, tx, &user).
			Return(fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		err = srv.AddUser(ctx, &user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; rollback")
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

		user := entity.User{
			Name:     name,
			Password: password,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddUser(ctx, tx, &user).
			Return(fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		err = srv.AddUser(ctx, &user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; rollbackerr; rollback")
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

		user := entity.User{
			Name:     name,
			Password: password,
		}

		m.EXPECT().Tx().Return(tx, nil)
		m.
			EXPECT().
			AddUser(ctx, tx, &user).
			Return(nil)

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

		err = srv.AddUser(ctx, &user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not add user; commit failed")
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

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
		assert.Equal(t, "could not delete user; rollbackfail; deletefail", err.Error())
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

		errDeleteFail := errors.New("deletefail")
		m.
			EXPECT().
			DeleteEmailsByUserID(ctx, tx, userID).
			Return(errDeleteFail)

		mdb.ExpectRollback().WillReturnError(errors.New("rollbackfail"))

		err = srv.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, errDeleteFail))
		assert.Equal(t, "could not delete user emails; rollbackfail; deletefail", err.Error())
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

	m := mock.NewMockInterface(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "name"
	filter := store.FilterUsers{Limit: 10}
	ctx := context.Background()

	// success
	{
		var IDs []int64
		m.
			EXPECT().
			FilterUsersID(ctx, filter, &IDs).
			DoAndReturn(func(_ context.Context, _ store.FilterUsers, ids *[]int64) error {
				*ids = append(*ids, userID)
				return nil
			})
		m.
			EXPECT().
			FetchUsers(ctx, []int64{userID}, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []int64, users *[]entity.User) error {
				*users = append(*users, entity.User{
					ID:   userID,
					Name: name,
				})
				return nil
			})

		var users []entity.User
		err := srv.FilterUsers(ctx, filter, &users)
		assert.Nil(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, users[0].ID, userID)
		assert.Equal(t, users[0].Name, name)
	}

	// fail
	{
		var IDs []int64
		m.
			EXPECT().
			FilterUsersID(ctx, filter, &IDs).
			Return(fmt.Errorf("opz"))

		var users []entity.User
		err := srv.FilterUsers(ctx, filter, &users)
		assert.Len(t, users, 0)
		assert.Equal(t, err.Error(), "opz")
	}
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "userName"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FetchUsers(ctx, []int64{userID}, gomock.Any()).
			DoAndReturn(func(_ context.Context, ids []int64, users *[]entity.User) error {
				*users = []entity.User{
					{
						ID:   userID,
						Name: name,
					},
				}
				return nil
			})

		var user entity.User
		err := srv.GetUserByID(ctx, userID, &user)
		assert.Nil(t, err)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, userID, user.ID)
	}

	// fails if storage fails
	{
		m.
			EXPECT().
			FetchUsers(ctx, []int64{userID}, gomock.Any()).
			Return(fmt.Errorf("opz"))

		var user entity.User
		err := srv.GetUserByID(ctx, userID, &user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}

	// fails if no user found
	{
		m.
			EXPECT().
			FetchUsers(ctx, []int64{userID}, gomock.Any()).
			Return(nil)

		var user entity.User
		err := srv.GetUserByID(ctx, userID, &user)
		fmt.Println("?", err)
		assert.True(t, errors.Is(err, errors.ErrNotFound))
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockInterface(ctrl)

	srv := service.New(&config.Config{}, m, nil)

	var userID int64 = 99
	name := "userName"
	email := "contact@example.com"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FilterUsersID(ctx, store.FilterUsers{Email: email, Limit: service.FilterUsersDefaultLimit}, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ store.FilterUsers, ids *[]int64) error {
				*ids = append(*ids, userID)
				return nil
			})
		m.
			EXPECT().
			FetchUsers(ctx, []int64{userID}, gomock.Any()).
			DoAndReturn(func(_ context.Context, ids []int64, users *[]entity.User) error {
				*users = append(*users, entity.User{
					ID:   userID,
					Name: name,
				})
				return nil
			})

		var user entity.User
		err := srv.GetUserByEmail(ctx, email, &user)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, name, user.Name)
	}

	// fails if storage fails
	{
		m.
			EXPECT().
			FilterUsersID(ctx, store.FilterUsers{Email: email, Limit: service.FilterUsersDefaultLimit}, gomock.Any()).
			Return(fmt.Errorf("opz"))

		var user entity.User
		err := srv.GetUserByEmail(ctx, email, &user)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz")
	}

	// fails if no user found
	{
		m.
			EXPECT().
			FilterUsersID(ctx, store.FilterUsers{Email: email, Limit: service.FilterUsersDefaultLimit}, gomock.Any()).
			Return(nil)

		var user entity.User
		err := srv.GetUserByEmail(ctx, email, &user)
		assert.Equal(t, errors.ErrNotFound, err)
	}
}
