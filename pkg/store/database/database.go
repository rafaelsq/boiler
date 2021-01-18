package database

import (
	"context"
	"database/sql"
	"fmt"

	"boiler/pkg/errors"

	"github.com/mattn/go-sqlite3"
)

// Database map a database access
type Database struct {
	sql *sql.DB
}

// Tx start a new transaction
func (s *Database) Tx() (*sql.Tx, error) {
	return s.sql.Begin()
}

// New create a new database
func New(sql *sql.DB) *Database {
	return &Database{
		sql: sql,
	}
}

// Insert execute an insert sql statement
func Insert(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if e, is := err.(sqlite3.Error); is && e.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, errors.ErrAlreadyExists
		}

		return 0, fmt.Errorf("could not insert; %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fail to retrieve last inserted ID; %w", err)
	}

	return id, nil
}

// Delete execute a delete sql statement
func Delete(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) error {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("could not remove; %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not fetch rows affected; %w", err)
	}

	if n == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Select execute a select sql statement
func Select(ctx context.Context, sql *sql.DB, scan func(func(...interface{}) error) (interface{}, error),
	query string, args ...interface{}) ([]interface{}, error) {

	rawRows, err := sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch rows; %w", err)
	}

	var rows []interface{}
	for {
		if !rawRows.Next() {
			break
		}

		row, err := scan(rawRows.Scan)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func scanInt(sc func(dest ...interface{}) error) (interface{}, error) {
	var id int64

	err := sc(&id)
	if err != nil {
		return nil, fmt.Errorf("could not scan int; %w", err)
	}

	return id, nil
}
