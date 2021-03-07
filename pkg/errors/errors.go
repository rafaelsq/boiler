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

	ErrBadRequest   = AddCodeWithMessage(nil, "bad_request", "bad request")
	ErrUnauthorized = AddCodeWithMessage(nil, "unauthorized", "unauthorized")

	// Service

	ErrNotFound            = AddCodeWithMessage(ErrBadRequest, "not_found", "not found")
	ErrAlreadyExists       = AddCodeWithMessage(ErrBadRequest, "already_exists", "already exists")
	ErrInvalidID           = AddCodeWithMessage(ErrBadRequest, "invalid_id", "invalid ID")
	ErrInvalidUserID       = AddCodeWithMessage(ErrInvalidID, "invalid_user_id", "invalid user ID")
	ErrInvalidEmailID      = AddCodeWithMessage(ErrInvalidID, "invalid_email_id", "invalid email ID")
	ErrInvalidPassword     = AddCodeWithMessage(ErrBadRequest, "invalid_password", "invalid password")
	ErrInvalidEmailAddress = AddCodeWithMessage(ErrBadRequest, "invalid_email_address", "invalid email address")
	ErrInvalidName         = AddCodeWithMessage(ErrBadRequest, "invalid_name", "invalid name")
	ErrInvalidLimit        = AddCodeWithMessage(ErrBadRequest, "invalid_limit", "invalid limit")
)
