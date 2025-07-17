package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already registered")
	ErrInvalidPassword   = errors.New("invalid password")
)
