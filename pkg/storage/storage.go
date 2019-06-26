package storage

import (
	"database/sql"

	"github.com/rafaelsq/boiler/pkg/iface"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	sql *sql.DB
}

func (s *Storage) Tx() (*sql.Tx, error) {
	return s.sql.Begin()
}

func New(sql *sql.DB) iface.Storage {
	return &Storage{
		sql: sql,
	}
}
