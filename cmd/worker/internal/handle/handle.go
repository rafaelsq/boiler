package handle

import (
	"boiler/pkg/service"
	"context"

	"github.com/gocraft/work"
)

func New(srv service.Interface) Handle {
	return Handle{srv}
}

type Handle struct {
	service service.Interface
}

func (h *Handle) DeleteUser(j *work.Job) error {
	return h.service.DeleteUser(context.Background(), j.ArgInt64("id"))
}

func (h *Handle) DeleteEmail(j *work.Job) error {
	return h.service.DeleteEmail(context.Background(), j.ArgInt64("id"))
}
