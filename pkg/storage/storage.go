package storage

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

var (
	storage iface.Storage

	users  []*entity.User
	emails []*entity.Email
)

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

	storage = &_storage{}
}

func GetDB() iface.Storage {
	return storage
}

type _storage struct{}

func (*_storage) AddEmail(ctx context.Context, userID int, address string) (int, error) {
	pk := int(len(emails))
	emails = append(emails, &entity.Email{ID: pk, UserID: userID, Address: address})
	return pk, nil
}

func (*_storage) Users(ctx context.Context) ([]*entity.User, error) {
	return users, nil
}

func (*_storage) UserByID(ctx context.Context, userID int) (*entity.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, nil
}

func (s *_storage) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, e := range emails {
		if e.Address == email {
			return s.UserByID(ctx, e.UserID)
		}
	}

	return nil, nil
}

func (*_storage) EmailsByUserID(ctx context.Context, userID int) ([]*entity.Email, error) {
	ret := []*entity.Email{}
	for _, email := range emails {
		if email.UserID == userID {
			ret = append(ret, email)
		}
	}

	return ret, nil
}
