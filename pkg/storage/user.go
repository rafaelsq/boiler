package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

func (s *Storage) AddUser(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	return Insert(ctx, tx, "INSERT INTO users (name, created, updated) VALUES (?, NOW(), NOW())", name)
}

func (s *Storage) DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	return Delete(ctx, tx, "DELETE FROM users WHERE id = ?", userID)
}

func (s *Storage) FilterUsersID(ctx context.Context, filter iface.FilterUsers) ([]int64, error) {
	limit := iface.FilterUsersDefaultLimit
	if filter.Limit != 0 {
		limit = filter.Limit
	}

	var args []interface{}
	var query string

	if len(filter.Email) != 0 {
		query = "SELECT u.id FROM users u INNER JOIN emails e ON(e.user_id = u.id) WHERE e.address = ?"
		args = append(args, filter.Email)
	} else {
		query = "SELECT id FROM users LIMIT ?"
		args = append(args, limit)
	}

	rows, err := Select(ctx, s.sql, scanInt, query, args...)
	if err != nil {
		return nil, err
	}

	IDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row != nil {
			IDs = append(IDs, row.(int64))
		}
	}

	return IDs, nil
}

func (s *Storage) FetchUsers(ctx context.Context, IDs ...int64) ([]*entity.User, error) {
	query := fmt.Sprintf(
		"SELECT id, name, UNIX_TIMESTAMP(created), UNIX_TIMESTAMP(updated) "+
			"FROM users WHERE id IN (%s) ORDER BY FIELD(id, %s)",
		strings.Repeat("?,", len(IDs))[0:len(IDs)*2-1],
		strings.Repeat("?,", len(IDs))[0:len(IDs)*2-1])

	args := make([]interface{}, 0, len(IDs)*2)
	for _, ID := range append(IDs, IDs...) {
		args = append(args, ID)
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
	var id int64
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
