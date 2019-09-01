package entity

import "time"

type Email struct {
	ID      int64     `json:"id"`
	UserID  int64     `json:"userID"`
	Address string    `json:"address"`
	Created time.Time `json:"created"`
}
