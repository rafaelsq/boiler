package user

import (
	"github.com/rafaelsq/boiler/pkg/entity"
)

type Repo struct {
	// db
}

func (r *Repo) ByID(id int) (*entity.User, error) {
	return &entity.User{
		ID: 1, Name: "John Doe",
	}, nil
}

func (r *Repo) FilterFriends(filter *entity.UserFriendsFilter) ([]*entity.User, error) {
	return []*entity.User{
		&entity.User{ID: 2, Name: "Little Finger Doe"},
		&entity.User{ID: 3, Name: "Aria Doe"},
	}, nil
}
