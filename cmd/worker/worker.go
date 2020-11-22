package main

import (
	"os"
	"os/signal"
	"time"

	"boiler/cmd"
	"boiler/cmd/worker/internal/handle"
	"boiler/pkg/service"
	"boiler/pkg/store/config"

	"github.com/gocraft/work"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.New()
	sv, redisPool := cmd.New(cfg)

	handler := handle.New(sv)

	pool := work.NewWorkerPool(handler, cfg.Worker.Concurrency, "all", redisPool)

	// middleware
	pool.Middleware(func(j *work.Job, next work.NextMiddlewareFunc) error {
		start := time.Now()
		log.Info().Str("job", j.Name).Msg("starting...")
		defer log.Info().Str("job", j.Name).Str("duration", time.Since(start).String()).Msg("finished")

		return next()
	})

	// Route
	pool.JobWithOptions(service.DeleteUser, work.JobOptions{Priority: 10, MaxFails: 1}, handler.DeleteUser)
	pool.JobWithOptions(service.DeleteEmail, work.JobOptions{Priority: 10, MaxFails: 1}, handler.DeleteEmail)

	// Start worker
	log.Info().Msg("[worker] Listening...")
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	pool.Stop()
}
