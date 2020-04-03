package api

import "github.com/milosgajdos/kraph/query"

const (
	// KindAll means all Kinds
	KindAll string = ""
	// NameAll means all names
	NameAll string = ""
	// NsALl means all namespaces
	NsAll string = ""
	// NamespaceNan means the resource is not namespaced
	NamespaceNan string = "nan"
)

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
	// Namespaced returns true if the resource is namespaced
	Namespaced() bool
}

// Relation defines remote link relation
type Relation interface {
	// String returns relation description
	String() string
}

// UID is object  UID
type UID interface {
	// String returns string UID
	String() string
}

// Link defines API object relation to another object
type Link interface {
	// To returns the UID of the object the link points to
	To() UID
	// Relation returns the type of the link relation
	Relation() Relation
}

// Object is an instance of a Resource
type Object interface {
	// UID is object uid
	UID() UID
	// Name is object name
	Name() string
	// Kind is Object kkind
	Kind() string
	// Namespace is object namespace
	Namespace() string
	// Links returns all object links
	Links() []Link
	// Raw returns a raw Object
	Raw() interface{}
}

// API is a map of all available API resources
type API interface {
	// Resources returns all API resources
	Resources() []Resource
	// Get returns all API resources matching the given query
	Get(...query.Option) ([]Resource, error)
}

// Top is an API topology i.e. the map of Objects
type Top interface {
	// Objects returns all objects in the topology
	Objects() []Object
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

// Client discovers API resources and maps API objects
type Client interface {
	Discoverer
	Mapper
}
