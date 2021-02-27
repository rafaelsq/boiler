package errors

import (
	"errors"
	"fmt"
)

// Wrap will add a error to a parent error.
func Wrap(parent, err error) error {
	return &WrapErr{Parent: parent, Err: err}
}

// WrapWithMessage will add a error with a message to a parent error.
func WrapWithMessage(parent, err error, message string) error {
	return &WrapErr{Parent: parent, Err: err, Msg: message}
}

// WrapErr is a error that wraps another error and support errors.[Is, As, etc]
// You can use it so that you can append a new error that can be unwrap through errors.Unwrap
type WrapErr struct {
	Parent error
	Err    error
	Msg    string
}

func (e *WrapErr) Unwrap() error {
	return e.Parent
}

func (e *WrapErr) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("%v %v; %v", e.Msg, e.Err, e.Parent)
	}
	return fmt.Sprintf("%v; %v", e.Err, e.Parent)
}

func (e *WrapErr) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *WrapErr) As(target interface{}) bool {
	if errors.As(e.Err, target) {
		return true
	}

	return errors.As(e.Parent, target)
}
