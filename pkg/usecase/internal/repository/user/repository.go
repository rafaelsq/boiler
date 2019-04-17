package user

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

type Repo struct {
	// db
}

func (r *Repo) ByID(ctx context.Context, id int) (*entity.User, error) {
	return &entity.User{
		ID: 1, Name: "John Doe",
	}, nil
}

func (r *Repo) FilterFriends(ctx context.Context, filter *entity.UserFriendsFilter) ([]*entity.User, error) {
	return []*entity.User{
		&entity.User{ID: 2, Name: "Little Finger Doe"},
		&entity.User{ID: 3, Name: "Aria Doe"},
	}, nil
}
