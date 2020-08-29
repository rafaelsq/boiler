// Package iface contains all the interface of the system
//go:generate ../../mock.sh
package iface

import (
	"context"
	"net/http"

	"github.com/rafaelsq/boiler/pkg/entity"
)

// Service is the interface of the Service
type Service interface {
	AddUser(context.Context, string, string) (int64, error)
	DeleteUser(context.Context, int64) error
	FilterUsers(context.Context, FilterUsers) ([]*entity.User, error)
	GetUserByID(context.Context, int64) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)
	AuthUser(context.Context, string, string) (*entity.User, string, error)

	FilterEmails(context.Context, FilterEmails) ([]*entity.Email, error)
	AddEmail(context.Context, int64, string) (int64, error)
	DeleteEmail(context.Context, int64) error

	AuthUserMiddleware(http.Handler) http.Handler
}
