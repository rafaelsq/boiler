package entity

import "context"

type Email struct {
	ID      uint   `json:"id"`
	Address string `json:"address"`
	User    User   `json:"user"`
}

type EmailUsecase interface {
	ByUserID(context.Context, uint) ([]*Email, error)
}

type EmailRepository interface {
	ByUserID(context.Context, uint) ([]*Email, error)
}
