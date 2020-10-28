package query

import "errors"

var (
	// ErrInvalidName is returned when name could not be decoded from query
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidGroup is returned when group could not be decoded from query
	ErrInvalidGroup = errors.New("invalid group")
	// ErrInvalidVersion is returned when version could not be decoded from query
	ErrInvalidVersion = errors.New("invalid version")
)
