package entity

import "context"

type User struct {
	ID   int
	Name string
}

type UserFilterOrder uint

var (
	UserFilterOrderASC  UserFilterOrder = 0
	UserFilterOrderDESC UserFilterOrder = 1
)

type UserFilter struct {
	Order UserFilterOrder
}

type UserUsecase interface {
	ByID(context.Context, int) (*User, error)
	Filter(context.Context, *UserFilter) ([]*User, error)
}

type UserRepository interface {
	ByID(context.Context, int) (*User, error)
	Filter(context.Context, *UserFilter) ([]*User, error)
}
