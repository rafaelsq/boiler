package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

func (s *Storage) AddUser(ctx context.Context, tx *sql.Tx, name string) (int, error) {
	result, err := tx.ExecContext(ctx,
		"INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())",
		name,
	)
	if err != nil {
		return 0, errors.New("could not insert user").SetArg("name", name).SetParent(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("last insert id failed after add user").SetParent(err)
	}

	return int(id), nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int) error {
	result, err := s.sql.ExecContext(ctx,
		"DELETE FROM users WHERE id = ?",
		userID,
	)
	if err != nil {
		return errors.New("could not remove user").SetArg("userID", userID).SetParent(err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not fetch rows affected after remove user").SetArg("userID", userID).SetParent(err)
	}

	if n == 0 {
		return errors.New("no rows affected").SetArg("userID", userID).SetParent(iface.ErrNotFound)
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
		return nil, errors.New("could not list users").SetArg("limit", limit).SetParent(err)
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
		return nil, errors.New("could not fetch user").SetArg("userID", userID).SetParent(err)
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
		return nil, errors.New("could not fetch user").SetArg("email", email).SetParent(err)
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
		return nil, errors.New("could not scan user").SetParent(err)
	}

	return &entity.User{
		ID:      id,
		Name:    name,
		Created: created,
		Updated: updated,
	}, nil
}
