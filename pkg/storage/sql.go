package storage

import (
	"database/sql"
	"time"

	// mariadb
	_ "github.com/go-sql-driver/mysql"
)

func NewMariaDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)
	db.SetConnMaxLifetime(time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}
