package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
)

// Resource is GitHub API resource
type Resource struct {
	*gen.Resource
}

// NewResource creates a new GitHub API resource and returns it
func NewResource(name, kind, group, version string, namespaced bool, opts api.Options) *Resource {
	return &Resource{
		Resource: gen.NewResource(name, kind, group, version, namespaced, opts),
	}
}
