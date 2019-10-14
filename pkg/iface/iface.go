package iface

const (
	// FilterUsersDefaultLimit is the default limit for user filtering
	FilterUsersDefaultLimit  uint = 50
	// FilterEmailsDefaultLimit is the default limit for email filtering 
	FilterEmailsDefaultLimit uint = 50
)

// FilterUsers is the input for filter users
type FilterUsers struct {
	Email  string
	Offset uint
	Limit  uint
}

// FilterEmails is the input for filter emails
type FilterEmails struct {
	EmailID int64
	UserID  int64
	Offset  uint
	Limit   uint
}
