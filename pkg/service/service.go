package service

import (
	"boiler/pkg/store"
	"boiler/pkg/store/config"

	"github.com/gocraft/work"
)

// New return a new service
func New(conf *config.Config, str store.Interface, enqueuer *work.Enqueuer) Interface {
	return &Service{
		enqueuer,
		conf,
		str,
	}
}

// Service is the main service
type Service struct {
	enqueuer *work.Enqueuer
	config   *config.Config
	store    store.Interface
}
