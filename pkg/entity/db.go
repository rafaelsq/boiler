package entity

import "context"

type DB interface {
	Users(context.Context) ([]*User, error)
	UserByID(context.Context, int) (*User, error)
	UserByEmail(context.Context, string) (*User, error)
	AddEmail(context.Context, int, string) (int, error)
	EmailsByUserID(context.Context, int) ([]*Email, error)
}
