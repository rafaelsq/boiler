package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/errors"
	"github.com/rafaelsq/boiler/pkg/iface"
	"go.uber.org/multierr"
)

func (s *Storage) AddUser(ctx context.Context, tx *sql.Tx, name string) (int, error) {
	result, err := tx.ExecContext(ctx,
		"INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())",
		name,
	)
	if err != nil {
		return 0, multierr.Append(err, errors.WithArg("could not insert user", "name", name))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, multierr.Append(err, errors.New("last insert id failed after add user"))
	}

	return int(id), nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int) error {
	result, err := s.sql.ExecContext(ctx,
		"DELETE FROM users WHERE id = ?",
		userID,
	)
	if err != nil {
		return multierr.Append(err, errors.WithArg("could not remove user", "userID", userID))
	}

	n, err := result.RowsAffected()
	if err != nil {
		return multierr.Append(err, errors.WithArg("could not fetch rows affected after remove user", "userID", userID))
	}

	if n == 0 {
		return multierr.Append(iface.ErrNotFound, errors.WithArg("no rows affected", "userID", userID))
	}

	return nil
}

func (s *Storage) FilterUsers(ctx context.Context, filter iface.FilterUsers) ([]*entity.User, error) {
	limit := iface.FilterUsersDefaultLimit
	if filter.Limit != 0 {
		limit = filter.Limit
	}

	if filter.UserID > 0 {
		u, err := s.filterUsersByID(ctx, filter.UserID)
		if err != nil || u == nil {
			return nil, err
		}

		return []*entity.User{u}, nil
	}

	if len(filter.Email) != 0 {
		u, err := s.filterUsersByEmail(ctx, filter.Email)
		if err != nil || u == nil {
			return nil, err
		}

		return []*entity.User{u}, nil
	}

	return s.filterUsers(ctx, limit)
}

func (s *Storage) filterUsers(ctx context.Context, limit uint) ([]*entity.User, error) {
	rows, err := s.sql.QueryContext(
		ctx,
		"SELECT id, name, created, updated FROM users LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, multierr.Append(err, errors.WithArg("could not list users", "limit", limit))
	}

	users := make([]*entity.User, 0, limit)
	for {
		if !rows.Next() {
			break
		}

		user, err := scanUser(rows.Scan)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) filterUsersByID(ctx context.Context, userID int) (*entity.User, error) {
	rows, err := s.sql.QueryContext(
		ctx,
		"SELECT id, name, created, updated FROM users WHERE id = ?",
		userID,
	)
	if err != nil {
		return nil, multierr.Append(err, errors.WithArg("could not fetch user", "userID", userID))
	}

	if !rows.Next() {
		return nil, nil
	}

	return scanUser(rows.Scan)
}

func (s *Storage) filterUsersByEmail(ctx context.Context, email string) (*entity.User, error) {
	rows, err := s.sql.QueryContext(
		ctx,
		"SELECT u.id, name, created, updated FROM users u INNER JOIN emails ON(user_id = u.id) WHERE email = ?",
		email,
	)
	if err != nil {
		return nil, multierr.Append(err, errors.WithArg("could not fetch user", "email", email))
	}

	if !rows.Next() {
		return nil, nil
	}

	return scanUser(rows.Scan)
}

func scanUser(sc func(dest ...interface{}) error) (*entity.User, error) {
	var id int
	var name string
	var created time.Time
	var updated time.Time

	err := sc(&id, &name, &created, &updated)
	if err != nil {
		return nil, multierr.Append(err, errors.New("could not scan user"))
	}

	return &entity.User{
		ID:      id,
		Name:    name,
		Created: created,
		Updated: updated,
	}, nil
}
