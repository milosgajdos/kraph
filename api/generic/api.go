package generic

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
)

// API is a generic API
type API struct {
	resources []api.Resource
	// resourceMap serves as an index into APIs
	resourceMap map[string][]api.Resource
}

// NewAPI returns new K8s API object
func NewAPI() *API {
	return &API{
		resources:   make([]api.Resource, 0),
		resourceMap: make(map[string][]api.Resource),
	}
}

// AddResource adds resource to the API
func (a *API) AddResource(r Resource) {
	a.resources = append(a.resources, r)
}

// AddResourceToPath adds resource to given path
func (a *API) AddResourceToPath(r Resource, path string) {
	a.resourceMap[path] = append(a.resourceMap[path], r)
}

// Resources returns all API resources
func (a *API) Resources() []api.Resource {
	resources := make([]api.Resource, len(a.resources))

	for i, res := range a.resources {
		resources[i] = res
	}

	return resources
}

// lookupGV searches all resources matching the given name and/or version
func (a *API) lookupGV(group, version string) ([]api.Resource, error) {
	var resources []api.Resource

	for _, r := range a.resources {
		if len(group) == 0 || group == r.Group() {
			if len(version) == 0 || version == r.Version() {
				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

// Get returns all API resources matching the given query
func (a *API) Get(opts ...query.Option) ([]api.Resource, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var apiResources []api.Resource

	if a.resourceMap == nil {
		return apiResources, nil
	}

	name := strings.ToLower(query.Name)
	group := query.Group
	version := query.Version

	if len(name) > 0 {
		if len(group) > 0 {
			if len(version) > 0 {
				return res2APIres(a.resourceMap[strings.Join([]string{name, group, version}, "/")])
			}
			return res2APIres(a.resourceMap[strings.Join([]string{name, group}, "/")])
		}

		// NOTE: both group and version are empty strings
		// we can safely return the resources indexed by name
		if len(version) == 0 {
			return res2APIres(a.resourceMap[name])
		}
	}

	return a.lookupGV(group, version)
}

// res2APIres "converts" resources into API resources
func res2APIres(rx []api.Resource) ([]api.Resource, error) {
	ax := make([]api.Resource, len(rx))

	for i, r := range rx {
		ax[i] = r
	}

	return ax, nil
}
