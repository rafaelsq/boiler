package database

import (
	"context"
	"database/sql"
	"time"

	"boiler/pkg/entity"
	"boiler/pkg/iface"

	"github.com/rafaelsq/errors"
)

// AddEmail insert a new emails in the database
func (s *Database) AddEmail(ctx context.Context, tx *sql.Tx, userID int64, address string) (int64, error) {
	return Insert(ctx, tx,
		"INSERT INTO emails (user_id, address, created) VALUES (?, ?, ?)",
		userID, address, time.Now(),
	)
}

// DeleteEmail remove an email from the database
func (s *Database) DeleteEmail(ctx context.Context, tx *sql.Tx, emailID int64) error {
	return Delete(ctx, tx, "DELETE FROM emails WHERE id = ?", emailID)
}

// DeleteEmailsByUserID remove email from the database
func (s *Database) DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	return Delete(ctx, tx, "DELETE FROM emails WHERE user_id = ?", userID)
}

// FilterEmails find for emails
func (s *Database) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	args := []interface{}{filter.UserID}
	where := "user_id = ?"
	if filter.EmailID > 0 {
		where = "id = ?"
		args = []interface{}{filter.EmailID}
	}

	rows, err := Select(ctx, s.sql, scanEmail,
		"SELECT id, user_id, address, created FROM emails WHERE "+where,
		args...,
	)
	if err != nil {
		return nil, err
	}

	emails := make([]*entity.Email, 0, len(rows))
	for _, row := range rows {
		emails = append(emails, row.(*entity.Email))

	}

	return emails, nil
}

func scanEmail(sc func(dest ...interface{}) error) (interface{}, error) {
	var id int64
	var userID int64
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
