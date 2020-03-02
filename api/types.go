package api

import "github.com/milosgajdos/kraph/query"

type Object interface {
	Raw() interface{}
}

type Resource interface {
	Name() string
	Group() string
	Version() string
	Namespaced() bool
}

type API interface {
	Resources() []Resource
	Lookup(string) []Resource
}

type Discoverer interface {
	Discover() (API, error)
}

type Mapper interface {
	Map(API) error
	Get(...query.Option) ([]Object, error)
}

type Client interface {
	Discoverer
	Mapper
}
