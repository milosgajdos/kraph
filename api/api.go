package api

import "github.com/milosgajdos/kraph/query"

// Object is an instance of a Resource
type Object interface {
	Name() string
	Kind() string
	Namespace() string
	Raw() interface{}
}

// Resource is an API resource
type Resource interface {
	Name() string
	Kind() string
	Group() string
	Version() string
	Namespaced() bool
}

// API allows to query resources
type API interface {
	Resources() []Resource
	Lookup(string) []Resource
}

// Discoverere discovers API
type Discoverer interface {
	Discover() (API, error)
}

// Mapper maps the API into topology
type Mapper interface {
	Map(API) (Top, error)
}

// Top is an API topology
type Top interface {
	Get(...query.Option) ([]Object, error)
	Raw() interface{}
}

// Client is API client
type Client interface {
	Discoverer
	Mapper
}
