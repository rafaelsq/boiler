// Package entity contains all the entities of the project
//go:generate msgp -tests=false
package entity

import "time"

// Email is the entity of the email
type Email struct {
	ID      int64     `json:"id"`
	UserID  int64     `json:"user_id"`
	Address string    `json:"address"`
	Created time.Time `json:"created"`
}
