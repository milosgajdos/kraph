package errors

import e "errors"

var (
	// ErrNotImplemented is returned by functions whose functionality has not been implemented yet
	ErrNotImplemented = e.New("not implemented")
	// ErrUnknownObject is returned when requesting an object which is not recognised
	ErrUnknownObject = e.New("unknown object")
)
