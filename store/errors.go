package store

import "errors"

var (
	// ErrNotFound is returned when store entity is not found
	ErrNotFound = errors.New("ErrNotFound")
)
