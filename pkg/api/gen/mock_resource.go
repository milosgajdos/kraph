package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
)

// NewMockResource returns new mock API resource
func NewMockResource(name, kind, group, version string, namespaced bool, opts api.Options) api.Resource {
	return NewResource(name, kind, group, version, namespaced, opts)
}
