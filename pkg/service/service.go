package service

import (
	"github.com/rafaelsq/boiler/pkg/iface"
)

// New return a new service
func New(storage iface.Storage) iface.Service {
	return &Service{storage}
}

// Service is the main service
type Service struct {
	storage iface.Storage
}
