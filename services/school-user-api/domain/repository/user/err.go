package user

import "errors"

var (
	// ErrDuplicateEntry happens when record violates unique constraints
	// defined in db
	ErrDuplicateEntry = errors.New("duplicate entry")
)
