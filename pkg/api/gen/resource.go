package gen

import "strings"

// Resource is a generic API resource
type Resource struct {
	name       string
	kind       string
	group      string
	version    string
	namespaced bool
}

// NewResrouce returns generic API resource
func NewResource(name, kind, group, version string, namespaced bool) *Resource {
	return &Resource{
		name:       name,
		kind:       kind,
		group:      group,
		version:    version,
		namespaced: namespaced,
	}
}

// Name returns the name of the resource
func (r Resource) Name() string {
	return r.name
}

// Kind returns resource kind
func (r Resource) Kind() string {
	return r.kind
}

// Group returns the API group of the resource
func (r Resource) Group() string {
	return r.group
}

// Version returns the version of the resource
func (r Resource) Version() string {
	return r.version
}

// Namespaced returns true if the resource is namespaced
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
