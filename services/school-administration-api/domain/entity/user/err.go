package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotValid    = errors.New("session not valid")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionHasExpired  = errors.New("session has expired")
)
