package storage_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestAddEmail(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	const address = "user@example.com"

	// succeed
	{
		userID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err = r.AddEmail(ctx, tx, userID, address)
		assert.Nil(t, err)
		assert.Equal(t, userID, 3)
		assert.Nil(t, tx.Commit())
	}

	// fail
	{
		userID := 3
		address := "user@example.com"
		myErr := fmt.Errorf("opz")

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnError(myErr)
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		emailID, err := r.AddEmail(ctx, tx, userID, address)
		assert.Equal(t, err.Error(), "could not insert email; opz")
		assert.Equal(t, emailID, 0)
		assert.Nil(t, tx.Commit())
	}

	// fails if duplicate
	{
		userID := 3
		address := "a@a.com"
		myErr := mysql.MySQLError{
			Message: "Duplicate entry 'a@a.com' for key 'address'",
			Number:  1062,
		}

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnError(&myErr)
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		emailID, err := r.AddEmail(ctx, tx, userID, address)
		assert.Equal(t, err, iface.ErrAlreadyExists)
		assert.Equal(t, emailID, 0)
		assert.Nil(t, tx.Commit())
	}

	// last insert failed
	{
		userID := 3
		address := "user@example.com"
		myErr := fmt.Errorf("opz")

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())"),
		).WithArgs(userID, address).WillReturnResult(sqlmock.NewResult(3, 1)).WillReturnResult(sqlmock.NewErrorResult(myErr))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		emailID, err := r.AddEmail(ctx, tx, userID, address)
		assert.Equal(t, err.Error(), "last insert id failed after add email address; opz")
		assert.Equal(t, emailID, 0)
		assert.Nil(t, tx.Commit())
	}
}

func TestDeleteEmail(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		emailID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmail(ctx, tx, emailID)
		assert.Nil(t, err)
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}

	// fails if exec fails
	{
		emailID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).WillReturnError(fmt.Errorf("opz"))

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmail(ctx, tx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not remove email; opz")
	}

	// fails if rows affected fails
	{
		emailID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("opz")))

		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmail(ctx, tx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not fetch rows affected after remove email; opz")
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}

	// fails if no rows affected
	{
		emailID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE id = ?"),
		).WithArgs(emailID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmail(ctx, tx, emailID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "no rows affected; not found")
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}
}

func TestDeleteEmailsByUserID(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		userID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE user_id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmailsByUserID(ctx, tx, userID)
		assert.Nil(t, err)
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}

	// fails
	{
		userID := 3

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM emails WHERE user_id = ?"),
		).WithArgs(userID).WillReturnError(fmt.Errorf("opz"))
		mock.ExpectCommit()

		r := storage.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteEmailsByUserID(ctx, tx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not remove emails by user ID; opz")
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}
}

func TestFilterEmails(t *testing.T) {
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

		r := storage.New(mdb)
		emails, err := r.FilterEmails(ctx, iface.FilterEmails{UserID: userID})
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

		r := storage.New(mdb)
		emails, err := r.FilterEmails(ctx, iface.FilterEmails{UserID: userID})
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

		r := storage.New(mdb)
		emails, err := r.FilterEmails(ctx, iface.FilterEmails{UserID: userID})
		assert.Equal(t, err.Error(), "could not fetch user's emails; opz")
		assert.Len(t, emails, 0)
	}
}
