package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/repository/email"
	"github.com/rafaelsq/boiler/pkg/repository/user"
)

func NewUser(db entity.DB) *User {
	return &User{
		db:        db,
		userRepo:  user.NewRepo(db),
		emailRepo: email.NewRepo(db),
	}
}

type User struct {
	db        entity.DB
	userRepo  entity.UserRepository
	emailRepo entity.EmailRepository
}

func (*User) ID(ctx context.Context, u *entity.User) (int, error) {
	return int(u.ID), nil
}

func (r *User) User(ctx context.Context, userID int) (*entity.User, error) {
	return r.userRepo.ByID(ctx, userID)
}

func (r *User) Users(ctx context.Context) ([]*entity.User, error) {
	return r.userRepo.Users(ctx)
}

func (r *User) Emails(ctx context.Context, u *entity.User) ([]*entity.Email, error) {
	return r.emailRepo.ByUserID(ctx, u.ID)
}
