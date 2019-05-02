package storage

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

var (
	db DB

	users  []*entity.User
	emails []*entity.Email
)

type DB interface {
	Users(context.Context) ([]*entity.User, error)
	UserByID(context.Context, int) (*entity.User, error)
	UserByEmail(context.Context, string) (*entity.User, error)
	AddEmail(context.Context, int, string) (int, error)
	EmailsByUserID(context.Context, int) ([]*entity.Email, error)
}

func init() {
	emails = []*entity.Email{
		&entity.Email{ID: 1, UserID: 1, Address: "john.doe@example.com"},
		&entity.Email{ID: 2, UserID: 1, Address: "johndoe@example.com"},
		&entity.Email{ID: 3, UserID: 3, Address: "aria@example.com"},
	}

	users = []*entity.User{
		&entity.User{ID: 1, Name: "John Doe"},
		&entity.User{ID: 2, Name: "Little Finger Doe"},
		&entity.User{ID: 3, Name: "Aria Doe"},
	}

	db = &_db{}
}

func GetDB() DB {
	return db
}

type _db struct{}

func (d *_db) AddEmail(ctx context.Context, userID int, address string) (int, error) {
	pk := int(len(emails))
	emails = append(emails, &entity.Email{ID: pk, UserID: userID, Address: address})
	return pk, nil
}

func (d *_db) Users(ctx context.Context) ([]*entity.User, error) {
	return users, nil
}

func (d *_db) UserByID(ctx context.Context, userID int) (*entity.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, nil
}

func (d *_db) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, e := range emails {
		if e.Address == email {
			return d.UserByID(ctx, e.UserID)
		}
	}

	return nil, nil
}

func (d *_db) EmailsByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	ret := []*entity.Email{}
	for _, email := range emails {
		if email.UserID == userID {
			ret = append(ret, email)
		}
	}

	return ret, nil
}
