package usecase

import (
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase/internal/repository/user"
)

func NewUser( /*db*/ ) interface {
	ByID(int) (*entity.User, error)
} {
	return &User{&user.Repo{ /*db*/ }}
}

type User struct {
	Repo interface {
		ByID(int) (*entity.User, error)
	}
}

func (u *User) ByID(id int) (*entity.User, error) {
	return u.Repo.ByID(id)
}
