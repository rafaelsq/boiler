package email_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/repository/email"
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
		userID := 3
		address := "user@example.com"

		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnResult(sqlmock.NewResult(3, 1))

		r := email.New(&StorageMock{mdb})
		userID, err := r.Add(ctx, userID, address)
		assert.Nil(t, err)
		assert.Equal(t, userID, 3)
	}

	// fail
	{
		userID := 3
		address := "user@example.com"
		myErr := fmt.Errorf("opz")

		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnError(myErr)

		r := email.New(&StorageMock{mdb})
		emailID, err := r.Add(ctx, userID, address)
		assert.Equal(t, err.Error(), "opz; could not insert email")
		assert.Equal(t, emailID, 0)
	}

	// fails if duplicate
	{
		userID := 3
		address := "a@a.com"
		myErr := mysql.MySQLError{
			Message: "Duplicate entry 'a@a.com' for key 'address'",
			Number:  1062,
		}

		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnError(&myErr)

		r := email.New(&StorageMock{mdb})
		emailID, err := r.Add(ctx, userID, address)
		assert.Equal(t, err, iface.ErrAlreadyExists)
		assert.Equal(t, emailID, 0)
	}

	// last insert failed
	{
		userID := 3
		address := "user@example.com"
		myErr := fmt.Errorf("opz")

		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnResult(sqlmock.NewResult(3, 1)).WillReturnResult(sqlmock.NewErrorResult(myErr))

		r := email.New(&StorageMock{mdb})
		emailID, err := r.Add(ctx, userID, address)
		assert.Equal(t, err.Error(), "opz; last insert id failed after add email address")
		assert.Equal(t, emailID, 0)
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
		emailID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(3, 1))

		r := email.New(&StorageMock{mdb})
		err := r.Delete(ctx, emailID)
		assert.Nil(t, err)
	}

	// fails if query fails
	{
		emailID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnError(fmt.Errorf("opz"))

		r := email.New(&StorageMock{mdb})
		err := r.Delete(ctx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz; could not remove email")
	}

	// fails if rows affect fails
	{
		emailID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("opz")))

		r := email.New(&StorageMock{mdb})
		err := r.Delete(ctx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "opz; could not fetch rows affected after remove email")
	}

	// fails if no rows affect
	{
		emailID := 3

		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(0, 0))

		r := email.New(&StorageMock{mdb})
		err := r.Delete(ctx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "not found; no rows affected")
	}
}

func TestByUserID(t *testing.T) {
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
			regexp.QuoteMeta("SELECT id, user_id, address, created FROM emails WHERE user_id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "address", "created"}).
				AddRow(3, userID, "user@example.com", time.Time{}),
		)

		r := email.New(&StorageMock{mdb})
		emails, err := r.ByUserID(ctx, userID)
		assert.Nil(t, err)
		assert.Len(t, emails, 1)
	}

	// scan fail
	{
		userID := 3

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, user_id, address, created FROM emails WHERE user_id = ?"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "address", "created"}).
				AddRow("opz", userID, "user@example.com", time.Time{}),
		)

		r := email.New(&StorageMock{mdb})
		emails, err := r.ByUserID(ctx, userID)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Len(t, emails, 0)
	}

	// fail
	{
		userID := 3
		myErr := fmt.Errorf("opz")

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id, user_id, address, created FROM emails WHERE user_id = ?"),
		).WithArgs(userID).WillReturnError(myErr)

		r := email.New(&StorageMock{mdb})
		emails, err := r.ByUserID(ctx, userID)
		assert.Equal(t, err.Error(), "opz; could not fetch user's emails")
		assert.Len(t, emails, 0)
	}
}
