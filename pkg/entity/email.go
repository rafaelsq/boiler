package entity

import "context"

type Email struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

type EmailService interface {
	ByUserID(context.Context, int) ([]*Email, error)
	Add(context.Context, int, string) (int, error)
}

type EmailRepository interface {
	ByUserID(context.Context, int) ([]*Email, error)
	Add(context.Context, int, string) (int, error)
}
