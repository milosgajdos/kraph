package generic

import (
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/metadata"
)

// Resource is a generic API resource
type Resource struct {
	name       string
	group      string
	version    string
	kind       string
	namespaced bool
	opts       api.Options
}

// NewResource creates a new generic API resource and returns it.
func NewResource(name, group, version, kind string, namespaced bool, opts api.Options) *Resource {
	return &Resource{
		name:       name,
		group:      group,
		version:    version,
		kind:       kind,
		namespaced: namespaced,
		opts:       opts,
	}
}

// Name returns the name of the resource
func (r Resource) Name() string {
	return r.name
}

// Group returns the API group of the resource
func (r Resource) Group() string {
	return r.group
}

// Version returns the version of the resource
func (r Resource) Version() string {
	return r.version
}

// Kind returns the resource kind
func (r Resource) Kind() string {
	return r.kind
}

// Namespaced returns true if the resource objects are namespaced
func (r Resource) Namespaced() bool {
	return r.namespaced
}

// Paths returns all possible variations of the resource paths
func (r Resource) Paths() []string {
	resNames := []string{strings.ToLower(r.name)}

	var names []string
	for _, name := range resNames {
		names = append(names,
			name,
			strings.Join([]string{name, r.group}, "/"),
			strings.Join([]string{name, r.group, r.version}, "/"),
		)
	}

	return names
}

// Metadata returns the resource metadata
func (r Resource) Metadata() metadata.Metadata {
	return r.opts.Metadata
}
