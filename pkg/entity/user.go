package entity

import "context"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	ByID(context.Context, int) (*User, error)
	ByEmail(context.Context, string) (*User, error)
	List(context.Context) ([]*User, error)
}

type UserRepository interface {
	ByID(context.Context, int) (*User, error)
	ByEmail(context.Context, string) (*User, error)
	List(context.Context) ([]*User, error)
}
