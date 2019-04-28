package entity

type Email struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	User    User   `json:"user"`
}
