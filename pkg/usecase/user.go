package usecase

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase/internal/repository/email"
	"github.com/rafaelsq/boiler/pkg/usecase/internal/repository/user"
)

func NewUser( /*db*/ ) entity.UserUsecase {
	return &User{
		UserRepo:  &user.Repo{ /*db*/ },
		EmailRepo: &email.Repo{ /*db*/ },
	}
}

type User struct {
	UserRepo  entity.UserRepository
	EmailRepo entity.EmailRepository
}

func (u *User) ByID(ctx context.Context, userID uint) (*entity.User, error) {
	done := make(chan error, 2)

	var user *entity.User
	go func() {
		var uErr error
		user, uErr = u.UserRepo.ByID(ctx, userID)
		if uErr != nil {
			done <- uErr
		}
		done <- nil
	}()

	var emails []*entity.Email
	go func() {
		var eErr error
		emails, eErr = u.EmailRepo.ByUserID(ctx, userID)
		if eErr != nil {
			done <- eErr
		}
		done <- nil
	}()

	err := <-done
	err2 := <-done
	if err != nil {
		return nil, err
	}

	if err2 != nil {
		return nil, err
	}

	user.Emails = emails
	for _, email := range emails {
		email.User = user
	}

	return user, nil
}

func (u *User) Filter(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return u.UserRepo.Filter(ctx, filter)
}
