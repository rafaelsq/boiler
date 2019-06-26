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

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	userID := 99
	name := "name"

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
			AddUser(ctx, tx, name).
			Return(userID, nil)
		mdb.ExpectCommit()

		id, err := srv.AddUser(ctx, name)
		assert.Nil(t, err)
		assert.Equal(t, userID, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}

	// fails if Tx fails
	{
		m.EXPECT().Tx().Return(nil, fmt.Errorf("opz"))

		id, err := srv.AddUser(ctx, name)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz; could not begin transaction")
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
			AddUser(ctx, tx, name).
			Return(0, fmt.Errorf("rollback"))
		mdb.ExpectRollback()

		id, err := srv.AddUser(ctx, name)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "rollback; could not add user")
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
			AddUser(ctx, tx, name).
			Return(0, fmt.Errorf("rollback"))

		mdb.ExpectRollback().WillReturnError(fmt.Errorf("rollbackerr"))

		id, err := srv.AddUser(ctx, name)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "rollback; rollbackerr; could not add user")
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
			AddUser(ctx, tx, name).
			Return(userID, nil)

		mdb.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

		id, err := srv.AddUser(ctx, name)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "commit failed; could not add user")
		assert.Equal(t, 0, id)
		assert.Nil(t, mdb.ExpectationsWereMet())
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	userID := 99

	ctx := context.Background()
	m.
		EXPECT().
		DeleteUser(ctx, userID).
		Return(nil)

	err := srv.DeleteUser(ctx, userID)
	assert.Nil(t, err)
}

func TestFilterUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	userID := 99
	name := "name"

	ctx := context.Background()
	m.
		EXPECT().
		FilterUsers(ctx, iface.FilterUsers{Limit: 10}).
		Return([]*entity.User{{
			ID:   userID,
			Name: name,
		}}, nil)

	vs, err := srv.FilterUsers(ctx, iface.FilterUsers{Limit: 10})
	assert.Nil(t, err)
	assert.Len(t, vs, 1)
	assert.Equal(t, vs[0].ID, userID)
	assert.Equal(t, vs[0].Name, name)
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	userID := 99
	name := "userName"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FilterUsers(ctx, iface.FilterUsers{UserID: userID}).
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
			FilterUsers(ctx, iface.FilterUsers{UserID: userID}).
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
			FilterUsers(ctx, iface.FilterUsers{UserID: userID}).
			Return([]*entity.User{}, nil)

		v, err := srv.GetUserByID(ctx, userID)
		assert.Nil(t, v)
		assert.Equal(t, iface.ErrNotFound, err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockStorage(ctrl)

	srv := service.New(m)

	userID := 99
	name := "userName"
	email := "contact@example.com"

	ctx := context.Background()

	// succeed
	{
		m.
			EXPECT().
			FilterUsers(ctx, iface.FilterUsers{Email: email}).
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
			FilterUsers(ctx, iface.FilterUsers{Email: email}).
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
			FilterUsers(ctx, iface.FilterUsers{Email: email}).
			Return([]*entity.User{}, nil)

		v, err := srv.GetUserByEmail(ctx, email)
		assert.Nil(t, v)
		assert.Equal(t, iface.ErrNotFound, err)
	}
}
