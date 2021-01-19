package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
)

// Resource is GitHub API resource
type Resource struct {
	*generic.Resource
}

// NewResource creates a new GitHub API resource and returns it
func NewResource(name, group, version, kind string, namespaced bool, opts api.Options) *Resource {
	return &Resource{
		Resource: generic.NewResource(name, group, version, kind, namespaced, opts),
	}
}
