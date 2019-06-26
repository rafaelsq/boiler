package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/errors"
	"github.com/rafaelsq/boiler/pkg/iface"
	"go.uber.org/multierr"
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

		return 0, multierr.Append(err, errors.WithArgs("could not insert email", map[string]interface{}{
			"userID":  userID,
			"address": address,
		}))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, multierr.Append(err, errors.New("last insert id failed after add email address"))
	}

	return int(id), nil
}

func (s *Storage) DeleteEmail(ctx context.Context, emailID int) error {
	result, err := s.sql.ExecContext(ctx,
		"DELETE FROM emails WHERE id = ?",
		emailID,
	)
	if err != nil {
		return multierr.Append(err, errors.WithArg("could not remove email", "emailID", emailID))
	}

	n, err := result.RowsAffected()
	if err != nil {
		return multierr.Append(err, errors.WithArg("could not fetch rows affected after remove email", "emailID", emailID))
	}

	if n == 0 {
		return multierr.Append(iface.ErrNotFound, errors.WithArg("no rows affected", "emailID", emailID))
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
		return nil, multierr.Append(err, errors.WithArg("could not fetch user's emails", "userID", filter.UserID))
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
		return nil, multierr.Append(err, errors.New("could not scan email"))
	}

	return &entity.Email{
		ID:      id,
		UserID:  userID,
		Address: address,
		Created: created,
	}, nil
}
