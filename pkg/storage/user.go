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
	return Insert(ctx, tx, "INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())", name)
}

func (s *Storage) DeleteUser(ctx context.Context, tx *sql.Tx, userID int) error {
	return Delete(ctx, tx, "DELETE FROM users WHERE id = ?", userID)
}

func (s *Storage) FilterUsers(ctx context.Context, filter iface.FilterUsers) ([]*entity.User, error) {
	limit := iface.FilterUsersDefaultLimit
	if filter.Limit != 0 {
		limit = filter.Limit
	}

	var args []interface{}
	var query string

	if filter.UserID > 0 {
		query = "SELECT id, name, created, updated FROM users WHERE id = ?"
		args = append(args, filter.UserID)
	} else if len(filter.Email) != 0 {
		query = "SELECT u.id, name, created, updated FROM users u INNER JOIN emails ON(user_id = u.id) WHERE email = ?"
		args = append(args, filter.Email)
	} else {
		query = "SELECT id, name, created, updated FROM users LIMIT ?"
		args = append(args, limit)
	}

	rows, err := Select(ctx, s.sql, scanUser, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, row.(*entity.User))

	}
	return users, nil
}

func scanUser(sc func(dest ...interface{}) error) (interface{}, error) {
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
