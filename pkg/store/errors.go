package store

import "errors"

var (
	// ErrUnknownEntity is returned when requesting an unknown entity
	ErrUnknownEntity = errors.New("unknown entity")
)
