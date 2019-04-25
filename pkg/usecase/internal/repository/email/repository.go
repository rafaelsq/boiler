package email

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
)

var emails []*entity.Email

func init() {
	emails = []*entity.Email{
		&entity.Email{ID: 1, User: &entity.User{ID: 1}, Address: "john.doe@example.com"},
		&entity.Email{ID: 2, User: &entity.User{ID: 1}, Address: "johndoe@example.com"},
		&entity.Email{ID: 3, User: &entity.User{ID: 3}, Address: "aria@example.com"},
	}
}

type Repo struct {
	// db
}

func (r *Repo) ByUserID(ctx context.Context, userID uint) ([]*entity.Email, error) {
	ret := []*entity.Email{}
	for _, email := range emails {
		if email.User.ID == userID {
			ret = append(ret, email)
		}
	}

	return ret, nil
}
