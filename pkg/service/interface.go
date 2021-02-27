//go:generate go run github.com/golang/mock/mockgen -package=mock -source=$GOFILE -destination=mock/service.go
package service

import (
	"context"
	"errors"

	"boiler/pkg/entity"
	"boiler/pkg/store"
)

var (
	// ErrInvalidID invalid ID error
	ErrInvalidID = errors.New("invalid ID")
	// ErrInvalidPassword invalid password
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	DeleteUser  string = "delete_user"
	DeleteEmail string = "delete_email"
)
const (
	// FilterUsersDefaultLimit is the default limit for user filtering
	FilterUsersDefaultLimit uint = 50
	// FilterEmailsDefaultLimit is the default limit for email filtering
	FilterEmailsDefaultLimit uint = 50
)

type Interface interface {
	AddUser(context.Context, *entity.User) error
	DeleteUser(context.Context, int64) error
	FilterUsers(context.Context, store.FilterUsers, *[]entity.User) error
	GetUserByID(context.Context, int64, *entity.User) error
	GetUserByEmail(context.Context, string, *entity.User) error
	AuthUser(context.Context, string, string, *entity.User, *string) error

	FilterEmails(context.Context, store.FilterEmails, *[]entity.Email) error
	AddEmail(context.Context, *entity.Email) error
	DeleteEmail(context.Context, int64) error
	EnqueueDeleteEmail(context.Context, int64) error
}
