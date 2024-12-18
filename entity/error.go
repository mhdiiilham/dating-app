package entity

import "errors"

// ErrInternalServerError var ...
var ErrInternalServerError = errors.New("INTERNAL_ERROR")

// Authentication related error.
var (
	ErrInvalidEmailAddress = errors.New("INVALID_EMAIL_ADDRESS")
	ErrInvalidAccessToken  = errors.New("INVALID_ACCESS_TOKEN")
	ErrUserAlreadyExist    = errors.New("USER_EXISTED")
)
