package entity

import "context"

type User struct {
	ID   int
	Name string
}

type UserFriendsFilter struct {
	FromUserID int
}

type UserUsecase interface {
	ByID(context.Context, int) (*User, error)
	Friends(context.Context, *UserFriendsFilter) ([]*User, error)
}

type UserRepository interface {
	ByID(context.Context, int) (*User, error)
	FilterFriends(context.Context, *UserFriendsFilter) ([]*User, error)
}
