package usecase

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase/internal/repository/user"
)

func NewUser( /*db*/ ) entity.UserUsecase {
	return &User{&user.Repo{ /*db*/ }}
}

type User struct {
	Repo entity.UserRepository
}

func (u *User) ByID(ctx context.Context, id int) (*entity.User, error) {
	return u.Repo.ByID(ctx, id)
}

func (u *User) Filter(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error) {
	return u.Repo.Filter(ctx, filter)
}
