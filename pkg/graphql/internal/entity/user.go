package entity

import "github.com/rafaelsq/boiler/pkg/entity"

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Emails []Email `json:"emails"`
}

func NewUser(u *entity.User) *User {
	return &User{
		ID:   u.ID,
		Name: u.Name,
	}
}
