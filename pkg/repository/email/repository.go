package email

import (
	"context"
	"time"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func New(storage iface.Storage) iface.EmailRepository {
	return &repository{storage}
}

type repository struct {
	storage iface.Storage
}

func (r *repository) Add(ctx context.Context, userID int, address string) (int, error) {
	result, err := r.storage.SQL().ExecContext(ctx,
		"INSERT INTO emails (user_id, address, created) VALUES (?, ?, NOW())",
		userID, address,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *repository) ByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	rows, err := r.storage.SQL().QueryContext(
		ctx,
		"SELECT id, user_id, address, created FROM emails WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}

	emails := make([]*entity.Email, 0)
	for {
		if !rows.Next() {
			break
		}

		email, err := scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

func scan(sc func(dest ...interface{}) error) (*entity.Email, error) {
	var id int
	var userID int
	var address string
	var created time.Time

	err := sc(&id, &userID, &address, &created)
	if err != nil {
		return nil, err
	}

	return &entity.Email{
		ID:      id,
		UserID:  userID,
		Address: address,
		Created: created,
	}, nil
}
