package entity

import "context"

type Email struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	User    User   `json:"user"`
}

type EmailRepository interface {
	ByUserID(context.Context, int) ([]*Email, error)
	Add(context.Context, int, string) (int, error)
}
