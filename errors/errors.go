package errors

import err "errors"

var (
	// ErrNotImplemented is returned when requesting functionality has not been implemented yet
	ErrNotImplemented = err.New("not implemented")
	// ErrEntityMissing is returned when entity is missing from query
	ErrEntityMissing = err.New("entity missing")
	// ErrUnknownObject is returned when requesting an unknown object
	ErrUnknownObject = err.New("unknown object")
	// ErrInvalidEntity is returned when requesting an invalid store entity
	ErrInvalidEntity = err.New("invalid entity")
	// ErrUnknownEntity is returned when requesting an unknown store entity
	ErrUnknownEntity = err.New("unknown entity")
	// ErrNodeNotFound is returned when a node could not be found
	ErrNodeNotFound = err.New("node not found")
	// ErrEdgeNotFound is returned when an edge could not be found
	ErrEdgeNotFound = err.New("edge not found")
	// ErrEdgeNotExist is returned when an edge does not exist
	ErrEdgeNotExist = err.New("edge does not exist")
	// ErrDuplicateNode is returned by store when duplicate nodes are found
	ErrDuplicateNode = err.New("duplicate node")
	// ErrMissingResource is returned by store when api.Object misses api.Resource
	ErrMissingResource = err.New("missing resource")
)
