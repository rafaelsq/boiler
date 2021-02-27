package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"boiler/pkg/entity"
	"boiler/pkg/store"
)

// AddEmail insert a new emails in the database
func (s *Database) AddEmail(ctx context.Context, tx *sql.Tx, email *entity.Email) error {
	id, err := Insert(ctx, tx,
		"INSERT INTO emails (user_id, address, created) VALUES (?, ?, ?)",
		email.UserID, email.Address, time.Now(),
	)
	email.ID = id
	return err
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
func (s *Database) FilterEmails(ctx context.Context, filter store.FilterEmails, emails *[]entity.Email) error {
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
		return err
	}

	*emails = make([]entity.Email, 0, len(rows))
	for _, row := range rows {
		*emails = append(*emails, *row.(*entity.Email))

	}

	return nil
}

func scanEmail(sc func(dest ...interface{}) error) (interface{}, error) {
	var id int64
	var userID int64
	var address string
	var created time.Time

	err := sc(&id, &userID, &address, &created)
	if err != nil {
		return nil, fmt.Errorf("could not scan email; %w", err)
	}

	return &entity.Email{
		ID:      id,
		UserID:  userID,
		Address: address,
		Created: created,
	}, nil
}
