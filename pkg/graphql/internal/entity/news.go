package entity

import "github.com/rafaelsq/boiler/pkg/entity"

func NewUser(u *entity.User) *User {
	return &User{
		ID:   u.ID,
		Name: u.Name,
	}
}

func NewEmail(e *entity.Email) *Email {
	return &Email{
		ID:      e.ID,
		Address: e.Address,
		User:    User{ID: e.UserID},
	}
}
