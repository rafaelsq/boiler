package user

import (
	"context"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/errors"
	"github.com/rafaelsq/boiler/pkg/iface"
	"go.uber.org/multierr"
)

func New(storage iface.Storage) iface.UserRepository {
	return &repository{storage}
}

type repository struct {
	storage iface.Storage
}

func (r *repository) Add(ctx context.Context, name string) (int, error) {
	result, err := r.storage.SQL().ExecContext(ctx,
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

func (r *repository) Delete(ctx context.Context, userID int) error {
	result, err := r.storage.SQL().ExecContext(ctx,
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

func (r *repository) List(ctx context.Context, limit uint) ([]*entity.User, error) {
	rows, err := r.storage.SQL().QueryContext(
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

		user, err := scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repository) ByID(ctx context.Context, userID int) (*entity.User, error) {
	rows, err := r.storage.SQL().QueryContext(
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

	return scan(rows.Scan)
}

func (r *repository) ByEmail(ctx context.Context, email string) (*entity.User, error) {
	rows, err := r.storage.SQL().QueryContext(
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

	return scan(rows.Scan)
}

func scan(sc func(dest ...interface{}) error) (*entity.User, error) {
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
