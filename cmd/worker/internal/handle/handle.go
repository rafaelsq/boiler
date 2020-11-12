package handle

import (
	"boiler/pkg/iface"
	"context"

	"github.com/gocraft/work"
)

func New(service iface.Service) Handle {
	return Handle{service}
}

type Handle struct {
	service iface.Service
}

func (h *Handle) DeleteUser(j *work.Job) error {
	return h.service.DeleteUser(context.Background(), j.ArgInt64("id"))
}

func (h *Handle) DeleteEmail(j *work.Job) error {
	return h.service.DeleteEmail(context.Background(), j.ArgInt64("id"))
}
