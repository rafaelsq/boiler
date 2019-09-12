package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// mariadb
	_ "github.com/go-sql-driver/mysql"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi"
	"github.com/rafaelsq/boiler/pkg/cache"
	"github.com/rafaelsq/boiler/pkg/router"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func newMariaDB(dsn string) (*sql.DB, error) {
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

func main() {
	var port = flag.Int("port", 2000, "")

	flag.Parse()

	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatal("Get RLIMIT_NOFILE failed", err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatal("Set RLIMIT_NOFILE failed", err)
	}

	mc := memcache.New("127.0.0.1:11211")

	sql, err := newMariaDB("root:boiler@tcp(127.0.0.1:3307)/boiler?timeout=5s&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	st := cache.New(mc, storage.New(sql))

	r := chi.NewRouter()
	router.ApplyMiddlewares(r)
	router.ApplyRoute(r, service.New(st))

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

		// start http shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("shutdown error", err)
		}

		close(iddleConnections)
	}()

	log.Printf("Listening on :%d\n", *port)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	log.Println("waiting iddle connections...")
	<-iddleConnections
	log.Println("done, bye!")
}
