package service

import (
	"boiler/pkg/iface"
	"boiler/pkg/store/config"

	"github.com/gocraft/work"
)

// New return a new service
func New(conf *config.Config, store iface.Store, enqueuer *work.Enqueuer) iface.Service {
	return &Service{
		enqueuer,
		conf,
		store,
	}
}

// Service is the main service
type Service struct {
	enqueue *work.Enqueuer
	config  *config.Config
	store   iface.Store
}
