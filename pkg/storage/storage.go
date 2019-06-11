package storage

import (
	"database/sql"
	"log"

	"github.com/rafaelsq/boiler/pkg/iface"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	sql *sql.DB
}

func (s *Storage) SQL() *sql.DB {
	return s.sql
}

func New(dsn string) iface.Storage {
	mariadb, err := NewMariaDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{
		sql: mariadb,
	}
}
