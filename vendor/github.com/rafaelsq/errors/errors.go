// Package errors provides simple error stack and a little more of context to the error
//
// Normally you would compose your error doing something like this
//
//     response, err := doSomething(somethingID)
//     if err != nil {
//         return fmt.Errorf(
//             "my error with id %d, response \"%v\"; %v",
//             somethingID,
//             response,
//             err,
//         )
//     }
//
// All the context from your error will have to be in the error string.
// After yout would need to extract all context from the string above.
//
// What this package do is add arguments to your error;
//
//     response, err := doSomething(somethingID)
//     if err != nil {
//         return errors.New("my error").
//              SetParent(err).
//              SetArg("somethingID", somethingID).
//              SetArg("response", response)
//     }
//
//     fmt.Println(err) // "something err; my error"
package errors

import "fmt"

// Error is the error with arguments
type Error struct {
	Msg    string
	Parent *error
	Args   map[string]interface{}
	Caller string
}

// Error returns erro as a string
func (e *Error) Error() string {
	if e.Parent != nil {
		return fmt.Sprintf("%v; %v", e.Msg, *e.Parent)
	}

	return e.Msg
}

// SetArg set an argument to the Error
func (e *Error) SetArg(arg string, value interface{}) *Error {
	e.Args[arg] = value
	return e
}

// SetParent set a parent error to Error
func (e *Error) SetParent(err error) *Error {
	e.Parent = &err
	return e
}

// Errorf formats according to a format specifier
func Errorf(format string, args ...interface{}) *Error {
	return &Error{
		Msg:    fmt.Sprintf(format, args...),
		Args:   make(map[string]interface{}),
		Caller: Caller(1),
	}
}

// New creates a new Error.
//
//     parentErr := New("parent")
//     err := New("error").SetArg("key", "value").SetParent(parentErr)
func New(err string) *Error {
	return &Error{
		Msg:    err,
		Args:   make(map[string]interface{}),
		Caller: Caller(1),
	}
}

// Cause returns the underlying cause of the error
func Cause(err error) error {
	if er, is := err.(*Error); is && er.Parent != nil {
		return Cause(*er.Parent)
	}

	return err
}

// List returns a list of errors
func List(err error) []error {
	if er, is := err.(*Error); is {
		if er.Parent != nil {
			return append(List(*er.Parent), er)
		}

		return []error{er}
	}

	return []error{err}
}
