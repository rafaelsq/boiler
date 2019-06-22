package user_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/stretchr/testify/assert"
)

type StorageMock struct{ sql *sql.DB }

func (s *StorageMock) SQL() *sql.DB {
	return s.sql
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		name := "user"

		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnResult(sqlmock.NewResult(3, 1))

		r := user.New(&StorageMock{mdb})
		userID, err := r.Add(ctx, name)
		assert.Nil(t, err)
		assert.Equal(t, userID, 3)
	}

	// fail
	{
		name := "user"

		myErr := fmt.Errorf("err")
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnError(myErr)

		r := user.New(&StorageMock{mdb})
		userID, err := r.Add(ctx, name)
		assert.Equal(t, err.Error(), "err; could not insert user")
		assert.Equal(t, userID, 0)
	}

	// last inserted failed
	{
		name := "user"

		myErr := fmt.Errorf("err")
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())"),
		).WithArgs(name).WillReturnResult(sqlmock.NewResult(3, 1)).WillReturnResult(sqlmock.NewErrorResult(myErr))

		r := user.New(&StorageMock{mdb})
		userID, err := r.Add(ctx, name)
		assert.Equal(t, err.Error(), "err; last insert id failed after add user")
		assert.Equal(t, userID, 0)
	}
}

func TestDelete(t *testing.T) {
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

		r := user.New(&StorageMock{mdb})
		err := r.Delete(ctx, userID)
		assert.Nil(t, err)
	}

	// fail if query fails
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnError(fmt.Errorf("opz"))

		r := user.New(&StorageMock{mdb})
		err := r.Delete(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz; could not remove user")
	}

	// fail if no rows affected
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 0))

		r := user.New(&StorageMock{mdb})
		err := r.Delete(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "not found; no rows affected")
	}

	// fail fatching rows affected
	{
		userID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("opz")))

		r := user.New(&StorageMock{mdb})
		err := r.Delete(ctx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz; could not fetch rows affected after remove user")
	}
}

func TestList(t *testing.T) {
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

		r := user.New(&StorageMock{mdb})
		users, err := r.List(ctx, limit)
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

		r := user.New(&StorageMock{mdb})
		users, err := r.List(ctx, limit)
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

		r := user.New(&StorageMock{mdb})
		users, err := r.List(ctx, limit)
		assert.Equal(t, err.Error(), "err; could not list users")
		assert.Len(t, users, 0)
	}
}

func TestByID(t *testing.T) {
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

		r := user.New(&StorageMock{mdb})
		user, err := r.ByID(ctx, userID)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, userID)
	}

	// succeed with no row
	{
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created", "updated"}),
		)

		r := user.New(&StorageMock{mdb})
		user, err := r.ByID(ctx, userID)
		assert.Nil(t, err)
		assert.Nil(t, user)
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

		r := user.New(&StorageMock{mdb})
		user, err := r.ByID(ctx, userID)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, user)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		userID := 3
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, name, created, updated FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnError(myErr)

		r := user.New(&StorageMock{mdb})
		user, err := r.ByID(ctx, userID)
		assert.Equal(t, err.Error(), "opz; could not fetch user")
		assert.Nil(t, user)
	}
}

func TestByMail(t *testing.T) {
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

		r := user.New(&StorageMock{mdb})
		user, err := r.ByEmail(ctx, email)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, 3)
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

		r := user.New(&StorageMock{mdb})
		user, err := r.ByEmail(ctx, email)
		assert.Nil(t, err)
		assert.Nil(t, user)
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

		r := user.New(&StorageMock{mdb})
		user, err := r.ByEmail(ctx, email)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, user)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id, name, created, updated FROM users u" +
				" INNER JOIN emails ON(user_id = u.id) WHERE email = ?"),
		).WithArgs(email).WillReturnError(myErr)

		r := user.New(&StorageMock{mdb})
		user, err := r.ByEmail(ctx, email)
		assert.Equal(t, err.Error(), "opz; could not fetch user")
		assert.Nil(t, user)
	}
}
