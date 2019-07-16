package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

func (s *Storage) AddEmail(ctx context.Context, tx *sql.Tx, userID int, address string) (int, error) {
	result, err := tx.ExecContext(ctx,
		"INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())",
		userID, address,
	)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				return 0, iface.ErrAlreadyExists
			}
		}

		return 0, errors.New("could not insert email").
			SetArg("userID", userID).
			SetArg("address", address).
			SetParent(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("last insert id failed after add email address").SetParent(err)
	}

	return int(id), nil
}

func (s *Storage) DeleteEmail(ctx context.Context, tx *sql.Tx, emailID int) error {
	result, err := tx.ExecContext(ctx,
		"DELETE FROM emails WHERE id = ?",
		emailID,
	)
	if err != nil {
		return errors.New("could not remove email").SetArg("emailID", emailID).SetParent(err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not fetch rows affected after remove email").
			SetArg("emailID", emailID).
			SetParent(err)
	}

	if n == 0 {
		return errors.New("no rows affected").SetArg("emailID", emailID).SetParent(iface.ErrNotFound)
	}

	return nil
}

func (s *Storage) DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int) error {
	_, err := tx.ExecContext(ctx,
		"DELETE FROM emails WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return errors.New("could not remove emails by user ID").SetArg("userID", userID).SetParent(err)
	}

	return nil
}

func (s *Storage) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	rows, err := s.sql.QueryContext(
		ctx,
		"SELECT id, user_id, address, created FROM emails WHERE user_id = ?",
		filter.UserID,
	)
	if err != nil {
		return nil, errors.New("could not fetch user's emails").SetArg("userID", filter.UserID).SetParent(err)
	}

	emails := make([]*entity.Email, 0)
	for {
		if !rows.Next() {
			break
		}

		email, err := scanEmail(rows.Scan)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

func scanEmail(sc func(dest ...interface{}) error) (*entity.Email, error) {
	var id int
	var userID int
	var address string
	var created time.Time

	err := sc(&id, &userID, &address, &created)
	if err != nil {
		return nil, errors.New("could not scan email").SetParent(err)
	}

	return &entity.Email{
		ID:      id,
		UserID:  userID,
		Address: address,
		Created: created,
	}, nil
}
