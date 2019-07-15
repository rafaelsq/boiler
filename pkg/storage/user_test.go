package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/storage"
	"github.com/stretchr/testify/assert"
)

type StorageMock struct{ sql *sql.DB }

func (s *StorageMock) SQL() *sql.DB {
	return s.sql
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		name := "user"

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name)
		assert.Nil(t, err)
		assert.Equal(t, userID, 3)
		assert.Nil(t, tx.Commit())
	}

	// fail
	{
		name := "user"

		myErr := fmt.Errorf("err")
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnError(myErr)
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name)
		assert.Equal(t, err.Error(), "could not insert user; err")
		assert.Equal(t, userID, 0)
		assert.Nil(t, tx.Commit())
	}

	// last inserted failed
	{
		name := "user"

		myErr := fmt.Errorf("err")
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnResult(sqlmock.NewResult(3, 1)).WillReturnResult(sqlmock.NewErrorResult(myErr))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name)
		assert.Equal(t, err.Error(), "last insert id failed after add user; err")
		assert.Equal(t, userID, 0)
		assert.Nil(t, tx.Commit())
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(3, 1))

		r := storage.New(mdb)
		err := r.DeleteUser(ctx, userID)
		assert.Nil(t, err)
	}

	// fail if query fails
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnError(fmt.Errorf("opz"))

		r := storage.New(mdb)
		err := r.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not remove user; opz")
	}

	// fail if no rows affected
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 0))

		r := storage.New(mdb)
		err := r.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "no rows affected; not found")
	}

	// fail fatching rows affected
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("opz")))

		r := storage.New(mdb)
		err := r.DeleteUser(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not fetch rows affected after remove user; opz")
	}
}

func TestFilterUsers(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		var limit uint = 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow(3, "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Limit: limit})
		assert.Nil(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, users[0].ID, 3)
	}

	// fail scan
	{
		var limit uint = 2
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow("err", "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Limit: limit})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Len(t, users, 0)
	}

	// fail
	{
		var limit uint = 4
		myErr := fmt.Errorf("err")

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnError(myErr)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Limit: limit})
		assert.Equal(t, err.Error(), "could not list users; err")
		assert.Len(t, users, 0)
	}
}

func TestFilterUsersByID(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow(userID, "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{UserID: userID})
		assert.Nil(t, err)
		assert.Equal(t, users[0].ID, userID)
	}

	// succeed with no row
	{
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{UserID: userID})
		assert.Nil(t, err)
		assert.Len(t, users, 0)
	}

	// scan fail
	{
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow("err", "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{UserID: userID})
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, users)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnError(myErr)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{UserID: userID})
		assert.Equal(t, err.Error(), "could not fetch user; opz")
		assert.Nil(t, users)
	}
}

func TestFilterUsersByMail(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id, name, created, updated FROM users u" +
				" INNER JOIN emails ON(user_id = u.id) WHERE email = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow(3, "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Email: email})
		assert.Nil(t, err)
		assert.Equal(t, users[0].ID, 3)
	}

	// succeed with no row
	{
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id, name, created, updated FROM users u" +
				" INNER JOIN emails ON(user_id = u.id) WHERE email = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Email: email})
		assert.Nil(t, err)
		assert.Len(t, users, 0)
	}

	// scan fail
	{
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id, name, created, updated FROM users u" +
				" INNER JOIN emails ON(user_id = u.id) WHERE email = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}).
				AddRow("err", "user", time.Time{}, time.Now()),
		)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Email: email})
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, users)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id, name, created, updated FROM users u" +
				" INNER JOIN emails ON(user_id = u.id) WHERE email = ?"),
		).WithArgs(email).WillReturnError(myErr)

		r := storage.New(mdb)
		users, err := r.FilterUsers(ctx, iface.FilterUsers{Email: email})
		assert.Equal(t, err.Error(), "could not fetch user; opz")
		assert.Nil(t, users)
	}
}
