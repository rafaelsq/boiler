package user

import (
	"context"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
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
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *repository) List(ctx context.Context, limit uint) ([]*entity.User, error) {
	rows, err := r.storage.SQL().QueryContext(
		ctx,
		"SELECT id, name, created, updated FROM users LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return &entity.User{
		ID:      id,
		Name:    name,
		Created: created,
		Updated: updated,
	}, nil
}
