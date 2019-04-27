package entity

import "context"

type User struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Emails []Email `json:"emails"`
}

type UserUsecase interface {
	ByID(context.Context, uint) (*User, error)
	ByEmail(context.Context, string) (*User, error)
	Users(context.Context) ([]*User, error)
}

type UserRepository interface {
	ByID(context.Context, uint) (*User, error)
	ByEmail(context.Context, string) (*User, error)
	Users(context.Context) ([]*User, error)
}
