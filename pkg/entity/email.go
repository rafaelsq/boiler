package entity

import "context"

type Email struct {
	ID      uint
	Address string
	User    *User
}

type EmailRepository interface {
	ByUserID(context.Context, uint) ([]*Email, error)
}
