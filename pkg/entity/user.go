// Package entity contains all the entities of the project
//go:generate go run github.com/tinylib/msgp -tests=false
package entity

import "time"

// User is the entity of the user
type User struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"-"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}
