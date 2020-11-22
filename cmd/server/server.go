package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"boiler/cmd"
	"boiler/cmd/server/internal/router"
	"boiler/pkg/store/config"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func main() {
	var port = flag.Int("port", 2000, "")

	flag.Parse()

	cfg := config.New()
	sv, _ := cmd.New(cfg)

	r := chi.NewRouter()
	router.ApplyMiddlewares(r, cfg, sv)
	router.ApplyRoute(r, sv)

	// graceful shutdown
	srv := http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: r}

	c := make(chan os.Signal, 1)
	iddleConnections := make(chan struct{})
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// sig is a ^C, handle it
		log.Warn().Msg("shutting down..")

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		go func() {
			<-c
			cancel()
		}()

		// start http shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("shutdown error")
		}

		close(iddleConnections)
	}()

	log.Info().Int("port", *port).Msg("[server] Listening...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Send()
	}

	log.Warn().Msg("waiting iddle connections...")
	<-iddleConnections
	log.Warn().Msg("done, bye!")
}
