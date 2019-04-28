package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/service"
)

func NewUser(db entity.DB) *User {
	return &User{
		db:           db,
		userService:  service.NewUser(db),
		emailService: service.NewEmail(db),
	}
}

type User struct {
	db           entity.DB
	userService  service.User
	emailService service.Email
}

func (*User) ID(ctx context.Context, u *entity.User) (int, error) {
	return int(u.ID), nil
}

func (r *User) User(ctx context.Context, userID int) (*entity.User, error) {
	return r.userService.ByID(ctx, userID)
}

func (r *User) Users(ctx context.Context) ([]*entity.User, error) {
	return r.userService.List(ctx)
}

func (r *User) Emails(ctx context.Context, u *entity.User) ([]*entity.Email, error) {
	return r.emailService.ByUserID(ctx, u.ID)
}
