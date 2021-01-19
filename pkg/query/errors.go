package query

import "errors"

var (
	// ErrInvalidUID is returned when uid could not be decoded from the query.
	ErrInvalidUID = errors.New("invalid uid")
	// ErrInvalidName is returned when name could not be decoded from the query.
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidGroup is returned when group could not be decoded from the query.
	ErrInvalidGroup = errors.New("invalid group")
	// ErrInvalidVersion is returned when version could not be decoded from query.
	ErrInvalidVersion = errors.New("invalid version")
	// ErrInvalidKind is returned when kind could not be decoded from query.
	ErrInvalidKind = errors.New("invalid kind")
	// ErrInvalidNamespace is returned when namespace could not be decoded from the query.
	ErrInvalidNamespace = errors.New("invalid namespace")
	// ErrInvalidEntity is returned when entity could not be decoded from the query.
	ErrInvalidEntity = errors.New("invalid entity")
)
