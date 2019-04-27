package boiler

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

var (
	db entity.DB

	users  []*entity.User
	emails []*entity.Email
)

func init() {
	emails = []*entity.Email{
		&entity.Email{ID: 1, User: entity.User{ID: 1}, Address: "john.doe@example.com"},
		&entity.Email{ID: 2, User: entity.User{ID: 1}, Address: "johndoe@example.com"},
		&entity.Email{ID: 3, User: entity.User{ID: 3}, Address: "aria@example.com"},
	}

	users = []*entity.User{
		&entity.User{ID: 1, Name: "John Doe"},
		&entity.User{ID: 2, Name: "Little Finger Doe"},
		&entity.User{ID: 3, Name: "Aria Doe"},
	}

	db = &DB{}
}

func GetDB() entity.DB {
	return db
}

type DB struct{}

func (d *DB) Users(ctx context.Context) ([]*entity.User, error) {
	return users, nil
}

func (d *DB) UserByID(ctx context.Context, userID uint) (*entity.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, nil
}

func (d *DB) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, e := range emails {
		if e.Address == email {
			return d.UserByID(ctx, e.User.ID)
		}
	}

	return nil, nil
}

func (d *DB) EmailsByUserID(ctx context.Context, userID uint) ([]*entity.Email, error) {
	ret := []*entity.Email{}
	for _, email := range emails {
		if email.User.ID == userID {
			ret = append(ret, email)
		}
	}

	return ret, nil
}
