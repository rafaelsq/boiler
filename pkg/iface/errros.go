package iface

import (
	"github.com/rafaelsq/errors"
)

var (
	ErrNotFound      = errors.New("not found").SetArg("code", "e0")
	ErrAlreadyExists = errors.New("already exists").SetArg("code", "s1")
)
