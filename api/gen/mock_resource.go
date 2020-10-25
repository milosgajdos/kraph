package gen

import "github.com/milosgajdos/kraph/api"

// NewMockResource returns new mock API resource
func NewMockResource(name, kind, group, version string, namespaced bool) api.Resource {
	return NewResource(name, kind, group, version, namespaced)
}
