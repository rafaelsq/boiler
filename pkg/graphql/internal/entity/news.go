package entity

import (
	"strconv"

	"github.com/rafaelsq/boiler/pkg/entity"
)

func NewUser(u *entity.User) *User {
	return &User{
		ID:   strconv.Itoa(u.ID),
		Name: u.Name,
	}
}

func NewEmail(e *entity.Email) *Email {
	return &Email{
		ID:      strconv.Itoa(e.ID),
		Address: e.Address,
		User:    &User{ID: strconv.Itoa(e.UserID)},
	}
}
