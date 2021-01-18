package database_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"boiler/pkg/errors"
	"boiler/pkg/store"
	"boiler/pkg/store/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
		password := "pass"

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, password, created, updated) VALUES (?, ?, ?, ?)"),
		).WithArgs(name, password, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name, password)
		assert.Nil(t, err)
		assert.Equal(t, 3, int(userID))
		assert.Nil(t, tx.Commit())
	}

	// fail
	{
		name := "user"
		password := "pass"

		myErr := fmt.Errorf("err")
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, password, created, updated) VALUES (?, ?, ?, ?)"),
		).WithArgs(name, password, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(myErr)
		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name, password)
		assert.Equal(t, err.Error(), "could not insert; err")
		assert.Equal(t, 0, int(userID))
		assert.Nil(t, tx.Commit())
	}

	// last inserted failed
	{
		name := "user"
		password := "pass"

		myErr := fmt.Errorf("err")
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO users (name, password, created, updated) VALUES (?, ?, ?, ?)"),
		).WithArgs(name, password, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(3, 1)).WillReturnResult(sqlmock.NewErrorResult(myErr))
		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		userID, err := r.AddUser(ctx, tx, name, password)
		assert.Equal(t, err.Error(), "fail to retrieve last inserted ID; err")
		assert.Equal(t, 0, int(userID))
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
		userID := int64(3)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnResult(sqlmock.NewResult(3, 1))
		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteUser(ctx, tx, userID)
		assert.Nil(t, err)
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}

	// fails if exec fails
	{
		userID := int64(3)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).WillReturnError(fmt.Errorf("opz"))

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteUser(ctx, tx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not remove; opz")
	}

	// fails if rows affected fails
	{
		userID := int64(3)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("opz")))

		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteUser(ctx, tx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "could not fetch rows affected; opz")
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}

	// fails if no rows affected
	{
		userID := int64(3)

		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta("DELETE FROM users WHERE id = ?"),
		).WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		r := database.New(mdb)

		tx, err := r.Tx()
		assert.Nil(t, err)

		err = r.DeleteUser(ctx, tx, userID)
		assert.NotNil(t, err)
		assert.Equal(t, err, errors.ErrNotFound)
		assert.Nil(t, tx.Commit())
		assert.Nil(t, mock.ExpectationsWereMet())
	}
}

func TestFilterUsersID(t *testing.T) {
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
			regexp.QuoteMeta("SELECT id FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(3),
		)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Limit: limit})
		assert.Nil(t, err)
		assert.Len(t, IDs, 1)
		assert.Equal(t, 3, int(IDs[0]))
	}

	// fail scan
	{
		var limit uint = 2
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow("err"),
		)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Limit: limit})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Len(t, IDs, 0)
	}

	// fail
	{
		var limit uint = 4
		myErr := fmt.Errorf("err")

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT id FROM users LIMIT ?"),
		).WithArgs(limit).WillReturnError(myErr)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Limit: limit})
		assert.Equal(t, "could not fetch rows; err", err.Error())
		assert.Len(t, IDs, 0)
	}
}

func TestFetchUsers(t *testing.T) {
	ctx := context.Background()
	mdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mdb.Close()

	// succeed
	{
		userID := int64(3)
		mock.ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, password, created, updated " +
					"FROM users WHERE id IN (?)"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "password", "created", "updated"}).
				AddRow(userID, "user", "pass", time.Time{}, time.Time{}),
		)

		r := database.New(mdb)
		users, err := r.FetchUsers(ctx, userID)
		assert.Nil(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, userID, users[0].ID)
		assert.Equal(t, "user", users[0].Name)
		assert.Equal(t, "pass", users[0].Password)
		assert.Equal(t, time.Time{}, users[0].Created)
		assert.Equal(t, time.Time{}, users[0].Updated)
	}

	// succeed with no row
	{
		userID := int64(3)
		mock.ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, password, created, updated " +
					"FROM users WHERE id IN (?)"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id"}),
		)

		r := database.New(mdb)
		users, err := r.FetchUsers(ctx, userID)
		assert.Nil(t, err)
		assert.Len(t, users, 0)
	}

	// empty
	{
		r := database.New(mdb)
		users, err := r.FetchUsers(ctx)
		assert.Nil(t, err)
		assert.Len(t, users, 0)
	}

	// scan fail
	{
		userID := int64(3)
		mock.ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, password, created, updated " +
					"FROM users WHERE id IN (?)"),
		).WithArgs(userID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "password", "created", "updated"}).
				AddRow("err", "user", "pass", 1, 2),
		)

		r := database.New(mdb)
		users, err := r.FetchUsers(ctx, userID)
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, users)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		userID := int64(3)
		mock.ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, password, created, updated " +
					"FROM users WHERE id IN (?)"),
		).WithArgs(userID).WillReturnError(myErr)

		r := database.New(mdb)
		users, err := r.FetchUsers(ctx, userID)
		assert.Equal(t, err.Error(), "could not fetch rows; opz")
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
			regexp.QuoteMeta("SELECT u.id FROM users u" +
				" INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(3),
		)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Email: email})
		assert.Nil(t, err)
		assert.Equal(t, 3, int(IDs[0]))
	}

	// succeed with no row
	{
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id FROM users u" +
				" INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id"}),
		)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Email: email})
		assert.Nil(t, err)
		assert.Len(t, IDs, 0)
	}

	// scan fail
	{
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id FROM users u" +
				" INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"),
		).WithArgs(email).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow("err"),
		)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Email: email})
		assert.Contains(t, err.Error(), "invalid syntax")
		assert.Nil(t, IDs)
	}

	// fail
	{
		myErr := fmt.Errorf("opz")
		email := "example@example.com"
		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT u.id FROM users u" +
				" INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"),
		).WithArgs(email).WillReturnError(myErr)

		r := database.New(mdb)
		IDs, err := r.FilterUsersID(ctx, store.FilterUsers{Email: email})
		assert.Equal(t, err.Error(), "could not fetch rows; opz")
		assert.Nil(t, IDs)
	}
}
