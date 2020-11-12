package cmd

import (
	"database/sql"
	"io/ioutil"
	"os"
	"time"

	"boiler/pkg/iface"
	"boiler/pkg/service"
	"boiler/pkg/store/config"
	"boiler/pkg/store/database"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func newDB(path string) (*sql.DB, error) {
	createTable := !fileExists(path)

	db, err := sql.Open("sqlite3", path)
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

	if createTable {
		log.Info().Msg("Creating DB")

		content, err := ioutil.ReadFile("pkg/store/database/schema.sql")
		if err != nil {
			db.Close()
			return nil, err
		}

		if _, err := db.Exec(string(content)); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func New() (iface.Service, *redis.Pool) {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := config.New()

	var redisPool = &redis.Pool{
		MaxActive: conf.Worker.Redis.MaxActive,
		MaxIdle:   conf.Worker.Redis.MaxIdle,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.Worker.Redis.Address)
		},
	}

	sql, err := newDB(conf.Sqlite3)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start DB")
	}

	st := database.New(sql)

	enqueuer := work.NewEnqueuer("all", redisPool)

	return service.New(conf, st, enqueuer), redisPool
}
