package entity

import "time"

type Email struct {
	ID      int       `json:"id"`
	UserID  int       `json:"userID"`
	Address string    `json:"address"`
	Created time.Time `json:"created"`
}
