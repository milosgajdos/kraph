package api

import "github.com/milosgajdos/kraph/query"

// Resource is an API resource
type Resource interface {
	// Name returns resource name
	Name() string
	// Kind returns resource kind
	Kind() string
	// Group retrurns resource group
	Group() string
	// Version returns resource version
	Version() string
	// Namespace returns resource namespace
	Namespaced() bool
}

// ObjRef is the reference object used when creating new object links
type ObjRef interface {
	// Name of the reference
	Name() string
	// Kind of the reference
	Kind() string
}

// Relation defines remote link relation
type Relation interface {
	// String returns relation description
	String() string
}

// Link defines API object relation to another object
type Link interface {
	// To returns the object reference the link points to
	To() ObjRef
	// Relation returns the type of link relation
	// NOTE: Relation is a widely used term in graph theory
	Relation() Relation
}

// Object is an instance of a Resource
type Object interface {
	// Name is object name
	Name() string
	// Kind is Object kkind
	Kind() string
	// Namespace is object namespace
	Namespace() string
	// Link links object to ObjRef
	Link(ObjRef, Relation) error
	// Links returns all object links
	Links() []Link
	// Raw returns a raw Objec that can be
	// typecasted into its Go type
	Raw() interface{}
}

// API allows to query API resources
type API interface {
	// Resources returns all API resources matching the given query
	Resources(...query.Option) []Resource
}

// Top is an API topology i.e. the map of Objects
type Top interface {
	// Get queries the topology and returns all matching objects
	Get(...query.Option) ([]Object, error)
}

// Discoverer discovers remote API
type Discoverer interface {
	// Discover returns the discovered API
	Discover() (API, error)
}

// Mapper maps the API into topology
type Mapper interface {
	// Map returns the API tpology
	Map(API) (Top, error)
}

// Client is API client
type Client interface {
	Discoverer
	Mapper
}
