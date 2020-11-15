package api

import (
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	// NsGlobal means global namespace
	NsGlobal string = "global"
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
	// Metadata returns Object metadata
	Metadata() metadata.Metadata
}

// Link defines API object relation to another object
type Link interface {
	// UID returns unique ID
	UID() uuid.UID
	// From returns the origin of the link
	From() uuid.UID
	// To returns the end the link links to
	To() uuid.UID
	// Metadata returns Object metadata
	Metadata() metadata.Metadata
}

// Object is an instance of a Resource
type Object interface {
	// UID is Object uniqque id
	UID() uuid.UID
	// Name is Object name
	Name() string
	// Namespace is Object namespace
	Namespace() string
	// Resource returns Object API resource
	Resource() Resource
	// Link links object to another object
	Link(uuid.UID, LinkOptions) error
	// Links returns all Object links
	Links() []Link
	// Metadata returns Object metadata
	Metadata() metadata.Metadata
}

// Source is the API source
type Source interface {
	// String returns API source as string
	String() string
}

// API is a map of all available API resources
type API interface {
	// Source is the API source
	Source() Source
	// Add adds resource to API
	Add(Resource, AddOptions) error
	// Resources returns all API resources
	Resources() []Resource
	// Get returns all API resources matching the query
	Get(*query.Query) ([]Resource, error)
}

// Top is an API topology i.e. the map of Objects
type Top interface {
	// API returns API source of topology
	API() API
	// Add adds Object to topology
	Add(Object, AddOptions) error
	// Objects returns all objects in the topology
	Objects() []Object
	// Get returns all API objects matching the query
	Get(*query.Query) ([]Object, error)
}

// Discoverer discovers an API
type Discoverer interface {
	// Discover discovers source and returns API
	Discover() (API, error)
}

// Mapper maps an API into topology
type Mapper interface {
	// Map returns the API tpology
	Map(API) (Top, error)
}

// Client discovers API resources and maps their objects
type Client interface {
	Discoverer
	Mapper
}
