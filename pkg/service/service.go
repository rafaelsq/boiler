package service

import (
	"boiler/pkg/iface"
	"boiler/pkg/store/config"
)

// New return a new service
func New(conf *config.Config, store iface.Store) iface.Service {
	return &Service{conf, store}
}

// Service is the main service
type Service struct {
	config *config.Config
	store  iface.Store
}
