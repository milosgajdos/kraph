package errors

import err "errors"

var (
	// ErrNotImplemented is returned when requesting functionality has not been implemented yet
	ErrNotImplemented = err.New("not implemented")
	// ErrUnknownObject is returned when requesting an unknown object
	ErrUnknownObject = err.New("unknown object")
	// ErrUnknownEntity is returned when requesting and unknown store entity
	ErrUnknownEntity = err.New("unknown entity")
	// ErrNodeNotFound is returned when a node could not be found
	ErrNodeNotFound = err.New("node not found")
	// ErrEdgeNotExist is returned when an edge does not exist
	ErrEdgeNotExist = err.New("edge does not exist")
	// ErrDuplicateNode is returned by store when duplicate nodes are found
	ErrDuplicateNode = err.New("duplicate node")
)
