package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"boiler/cmd/server/internal/router"
	"boiler/pkg/service"
	"boiler/pkg/store/config"
	"boiler/pkg/store/database"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
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
		log.Println("Creating DB")

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

func main() {
	var port = flag.Int("port", 2000, "")

	flag.Parse()

	sql, err := newDB("./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	st := database.New(sql)
	conf := config.New()

	sv := service.New(conf, st)

	r := chi.NewRouter()
	router.ApplyMiddlewares(r, sv)
	router.ApplyRoute(r, sv)

	// graceful shutdown
	srv := http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: r}

	c := make(chan os.Signal, 1)
	iddleConnections := make(chan struct{})
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// sig is a ^C, handle it
		log.Println("shutting down..")

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		go func() {
			<-c
			cancel()
		}()

		// start http shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("shutdown error", err)
		}

		close(iddleConnections)
	}()

	log.Printf("Listening on :%d\n", *port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("waiting iddle connections...")
	<-iddleConnections
	log.Println("done, bye!")
}
