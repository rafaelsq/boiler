package iface

import (
	"github.com/rafaelsq/errors"
)

var (
	// ErrNotFound not found error 
	ErrNotFound      = errors.New("not found").SetArg("code", "e0")
	// ErrAlreadyExists already exists error
	ErrAlreadyExists = errors.New("already exists").SetArg("code", "s1")
	// ErrInvalidID invalid ID error
	ErrInvalidID     = errors.New("invalid ID").SetArg("code", "iid")
)
