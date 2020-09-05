package entity

import (
	"strconv"

	"boiler/pkg/entity"
)

// NewUser return a new User entity
func NewUser(u *entity.User) *User {
	return &User{
		ID:   strconv.FormatInt(u.ID, 10),
		Name: u.Name,
	}
}

// NewEmail return a new Email entity
func NewEmail(e *entity.Email) *Email {
	return &Email{
		ID:      strconv.FormatInt(e.ID, 10),
		Address: e.Address,
		User:    &User{ID: strconv.FormatInt(e.UserID, 10)},
	}
}
