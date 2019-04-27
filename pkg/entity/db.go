package entity

import "context"

type DB interface {
	Users(context.Context) ([]*User, error)
	UserByID(context.Context, uint) (*User, error)
	UserByEmail(context.Context, string) (*User, error)
	EmailsByUserID(context.Context, uint) ([]*Email, error)
}
