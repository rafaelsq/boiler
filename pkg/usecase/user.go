package usecase

import (
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase/internal/repository/user"
)

func NewUser( /*db*/ ) entity.UserUsecase {
	return &User{&user.Repo{ /*db*/ }}
}

type User struct {
	Repo entity.UserRepository
}

func (u *User) ByID(id int) (*entity.User, error) {
	return u.Repo.ByID(id)
}

func (u *User) Friends(filter *entity.UserFriendsFilter) ([]*entity.User, error) {
	return u.Repo.FilterFriends(filter)
}
