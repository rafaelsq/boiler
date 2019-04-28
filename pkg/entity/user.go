package entity

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Emails []Email `json:"emails"`
}
