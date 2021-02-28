package errors

import (
	"errors"
)

var (
	As     = errors.As
	Is     = errors.Is
	Unwrap = errors.Unwrap
	New    = errors.New
)

var (
	// Base Errors

	ErrBadRequest   = AddCodeWithMessage(nil, "BAD_REQUEST", "bad request")
	ErrUnauthorized = AddCodeWithMessage(nil, "UNAUTHORIZED", "unauthorized")

	// Service

	ErrNotFound            = AddCodeWithMessage(ErrBadRequest, "NOT_FOUND", "not found")
	ErrAlreadyExists       = AddCodeWithMessage(ErrBadRequest, "ALREADY_EXISTS", "already exists")
	ErrInvalidID           = AddCodeWithMessage(ErrBadRequest, "INVALID_ID", "invalid ID")
	ErrInvalidPassword     = AddCodeWithMessage(ErrBadRequest, "INVALID_PASSWORD", "invalid password")
	ErrInvalidEmailAddress = AddCodeWithMessage(ErrBadRequest, "INVALID_EMAIL_ADDRESS", "invalid email address")
)
