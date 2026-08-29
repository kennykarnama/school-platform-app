package user

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotExist     = errors.New("user not exist")
)
