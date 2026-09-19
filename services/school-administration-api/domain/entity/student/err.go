package student

import "errors"

var (
	ErrAlternativeIDAlreadyExists   = errors.New("student alternative ID already exists")
	ErrActivePlacementAlreadyExists = errors.New("student already has an active placement for this academic year")
)
