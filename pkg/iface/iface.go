package iface

const (
	FilterUsersDefaultLimit  uint = 50
	FilterEmailsDefaultLimit uint = 50
)

type FilterUsers struct {
	UserID int64
	Email  string
	Offset uint
	Limit  uint
}

type FilterEmails struct {
	EmailID int64
	UserID  int64
	Offset  uint
	Limit   uint
}
