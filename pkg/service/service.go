package service

import (
	"github.com/rafaelsq/boiler/pkg/iface"
)

func New(storage iface.Storage) iface.Service {
	return &Service{storage}
}

type Service struct {
	storage iface.Storage
}
