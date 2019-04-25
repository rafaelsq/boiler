package user

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

var users []*entity.User

func init() {
	users = []*entity.User{
		&entity.User{ID: 1, Name: "John Doe"},
		&entity.User{ID: 2, Name: "Little Finger Doe"},
		&entity.User{ID: 3, Name: "Aria Doe"},
	}
}

type Repo struct {
	// db
}

func (r *Repo) ByID(ctx context.Context, id int) (*entity.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, nil
}

func (r *Repo) Filter(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return users, nil
}
