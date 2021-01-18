package errors

import (
	"errors"
	"fmt"
)

var (
	As     = errors.As
	Is     = errors.Is
	Unwrap = errors.Unwrap
	New    = errors.New
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrNotFound      = fmt.Errorf("not found; %w", ErrBadRequest)
	ErrAlreadyExists = fmt.Errorf("already exists; %w", ErrBadRequest)
)
