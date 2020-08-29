package service

import (
	"github.com/rafaelsq/boiler/pkg/config"
	"github.com/rafaelsq/boiler/pkg/iface"
)

// New return a new service
func New(conf *config.Config, storage iface.Storage) iface.Service {
	return &Service{conf, storage}
}

// Service is the main service
type Service struct {
	config  *config.Config
	storage iface.Storage
}
