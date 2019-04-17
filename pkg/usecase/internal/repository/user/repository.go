package user

import (
	"github.com/rafaelsq/boiler/pkg/entity"
)

type Repo struct {
	// db
}

func (*Repo) ByID(id int) (*entity.User, error) {
	return &entity.User{
		ID: 1, Name: "John Doe",
	}, nil
}
