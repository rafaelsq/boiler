package entity

import "github.com/rafaelsq/boiler/pkg/entity"

type Email struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	User    User   `json:"user"`
}

func NewEmail(e *entity.Email) *Email {
	return &Email{
		ID:      e.ID,
		Address: e.Address,
		User:    User{ID: e.UserID},
	}
}
