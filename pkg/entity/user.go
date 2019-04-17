package entity

type User struct {
	ID   int
	Name string
}

type UserFriendsFilter struct {
	FromUserID int
}

type UserUsecase interface {
	ByID(int) (*User, error)
	Friends(*UserFriendsFilter) ([]*User, error)
}

type UserRepository interface {
	ByID(int) (*User, error)
	FilterFriends(*UserFriendsFilter) ([]*User, error)
}
