package userauth

import "errors"

var (
	ErrNotAuthorized = errors.New("user is unauthorized")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenRevoked  = errors.New("token revoked")
)
