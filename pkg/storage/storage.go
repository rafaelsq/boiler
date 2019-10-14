package storage

import (
	"context"
	"database/sql"

	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // Preload mysql extension
)

// Storage map a database access
type Storage struct {
	sql *sql.DB
}

// Tx start a new transaction
func (s *Storage) Tx() (*sql.Tx, error) {
	return s.sql.Begin()
}

// New create a new storage
func New(sql *sql.DB) iface.Storage {
	return &Storage{
		sql: sql,
	}
}

// Insert execute an insert sql statement
func Insert(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				return 0, iface.ErrAlreadyExists
			}
		}

		return 0, errors.New("could not insert").SetArg("args", args).SetParent(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("fail to retrieve last inserted ID").SetParent(err)
	}

	return id, nil
}

// Delete execute a delete sql statement
func Delete(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) error {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.New("could not remove").SetArg("args", args).SetParent(err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not fetch rows affected").SetArg("args", args).SetParent(err)
	}

	if n == 0 {
		return iface.ErrNotFound
	}

	return nil
}

// Select execute a select sql statement
func Select(ctx context.Context, sql *sql.DB, scan func(func(...interface{}) error) (interface{}, error), query string, args ...interface{}) ([]interface{}, error) {
	rawRows, err := sql.QueryContext(ctx, query, args...)
	if err != nil {
		cerr := errors.New("could not fetch rows")
		cerr.Caller = errors.Caller(1)
		return nil, cerr.SetArg("args", args).SetParent(err)
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
		return nil, errors.New("could not scan int").SetParent(err)
	}

	return id, nil
}
