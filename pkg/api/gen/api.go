package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/query"
)

// API is a generic API
type API struct {
	// source is API source
	source string
	// resources are API resources
	resources []api.Resource
	// resourceMap serves as an index into APIs
	resourceMap map[string][]api.Resource
}

// NewAPI returns new K8s API object
func NewAPI(src string) *API {
	return &API{
		source:      src,
		resources:   make([]api.Resource, 0),
		resourceMap: make(map[string][]api.Resource),
	}
}

// AddResource adds resource to the API
func (a *API) AddResource(r api.Resource) {
	a.resources = append(a.resources, r)
}

// IndexPath indexes resource to given path
func (a *API) IndexPath(r api.Resource, path string) {
	a.resourceMap[path] = append(a.resourceMap[path], r)
}

// Source returns API source
func (a *API) Source() api.Source {
	return &Source{
		src: a.source,
	}
}

// Resources returns all API resources
func (a *API) Resources() []api.Resource {
	resources := make([]api.Resource, len(a.resources))

	for i, res := range a.resources {
		resources[i] = res
	}

	return resources
}

// Get returns all API resources matching the given query
func (a *API) Get(q *query.Query) ([]api.Resource, error) {
	var ar []api.Resource

	match := q.Matcher()

	for _, r := range a.resources {
		if match.NameVal(r.Name()) {
			if match.GroupVal(r.Group()) {
				if match.VersionVal(r.Version()) {
					ar = append(ar, r)
				}
			}
		}
	}

	return ar, nil
}
