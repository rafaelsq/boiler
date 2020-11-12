package main

import (
	"os"
	"os/signal"
	"time"

	"boiler/cmd"
	"boiler/cmd/worker/internal/handle"
	"boiler/pkg/iface"

	"github.com/gocraft/work"
	"github.com/rs/zerolog/log"
)

func main() {

	sv, redisPool := cmd.New()

	handler := handle.New(sv)

	pool := work.NewWorkerPool(handler, 10, "all", redisPool)

	// middleware
	pool.Middleware(func(j *work.Job, next work.NextMiddlewareFunc) error {
		start := time.Now()
		log.Info().Str("job", j.Name).Msg("starting...")
		defer log.Info().Str("job", j.Name).Str("duration", time.Since(start).String()).Msg("finished")

		return next()
	})

	// Route
	pool.JobWithOptions(iface.DeleteUser, work.JobOptions{Priority: 10, MaxFails: 1}, handler.DeleteUser)
	pool.JobWithOptions(iface.DeleteEmail, work.JobOptions{Priority: 10, MaxFails: 1}, handler.DeleteEmail)

	// Start worker
	log.Info().Msg("[worker] Listening...")
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	pool.Stop()
}
