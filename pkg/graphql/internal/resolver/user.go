package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func NewUser(us iface.UserService, es iface.EmailService) *User {
	return &User{
		userService:  us,
		emailService: es,
	}
}

type User struct {
	userService  iface.UserService
	emailService iface.EmailService
}

func (*User) ID(ctx context.Context, u *entity.User) (int, error) {
	return int(u.ID), nil
}

func (r *User) User(ctx context.Context, userID int) (*entity.User, error) {
	u, err := r.userService.ByID(ctx, userID)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, err
}

func (r *User) Users(ctx context.Context) ([]*entity.User, error) {
	us, err := r.userService.List(ctx)
	if err == nil {
		users := make([]*entity.User, 0, len(us))
		for _, u := range us {
			users = append(users, entity.NewUser(u))
		}
		return users, nil
	}
	return nil, err
}

func (r *User) Emails(ctx context.Context, u *entity.User) ([]*entity.Email, error) {
	es, err := r.emailService.ByUserID(ctx, u.ID)
	if err == nil {
		emails := make([]*entity.Email, 0, len(es))
		for _, e := range es {
			emails = append(emails, entity.NewEmail(e))
		}
		return emails, nil
	}
	return nil, err
}
