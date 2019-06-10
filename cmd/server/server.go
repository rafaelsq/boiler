package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	er "github.com/rafaelsq/boiler/pkg/repository/email"
	ur "github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/rafaelsq/boiler/pkg/router"
	"github.com/rafaelsq/boiler/pkg/service"
	"github.com/rafaelsq/boiler/pkg/storage"
)

func main() {
	var port = flag.Int("port", 2000, "")

	flag.Parse()

	us := service.NewUser(ur.New(storage.GetDB()))
	es := service.NewEmail(er.New(storage.GetDB()))

	r := chi.NewRouter()
	router.ApplyMiddlewares(r)
	router.ApplyRoute(r, us, es)

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
