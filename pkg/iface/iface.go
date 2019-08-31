package iface

const (
	FilterUsersDefaultLimit  uint = 50
	FilterEmailsDefaultLimit uint = 50
)

type FilterUsers struct {
	UserID int
	Email  string
	Offset uint
	Limit  uint
}

type FilterEmails struct {
	EmailID int
	UserID  int
	Offset  uint
	Limit   uint
}
