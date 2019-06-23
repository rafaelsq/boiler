package iface

import (
	"github.com/rafaelsq/boiler/pkg/errors"
)

var (
	ErrNotFound      = errors.WithArg("not found", "code", "e0")
	ErrAlreadyExists = errors.WithArg("already exists", "code", "s1")
)
