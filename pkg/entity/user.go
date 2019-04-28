package entity

import "context"

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Emails []Email `json:"emails"`
}

type UserRepository interface {
	ByID(context.Context, int) (*User, error)
	ByEmail(context.Context, string) (*User, error)
	Users(context.Context) ([]*User, error)
}
