package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
)

// API is API resource map
type API struct {
	resources []Resource
	// resourceMap serves as an index into APIs
	resourceMap map[string][]Resource
}

// Resources returns all API resources matching the given query
func (a *API) Resources(opts ...query.Option) []api.Resource {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var apiResources []api.Resource

	if a.resourceMap == nil {
		return apiResources
	}

	name := strings.ToLower(query.Name)
	group := query.Group
	version := query.Version

	// try pulling out the results from indexed entries
	if len(name) > 0 {
		if len(group) > 0 {
			if len(version) > 0 {
				return res2APIres(a.resourceMap[strings.Join([]string{name, group, version}, "/")])
			}
			return res2APIres(a.resourceMap[strings.Join([]string{name, group}, "/")])
		}
		// both group and version have 0 length; return the resources indexed by name
		if len(version) == 0 {
			return res2APIres(a.resourceMap[name])
		}
	}

	return a.lookupGV(group, version)
}

// lookupGV searches all resources matching the given name and/or version
func (a *API) lookupGV(group, version string) []api.Resource {
	var resources []api.Resource

	for _, r := range a.resources {
		if len(group) == 0 || group == r.gv.Group {
			if len(version) == 0 || version == r.gv.Version {
				resources = append(resources, r)
			}
		}
	}

	return resources
}

// res2APIres "converts" resources into API resources
func res2APIres(rx []Resource) []api.Resource {
	ax := make([]api.Resource, len(rx))

	for i, r := range rx {
		ax[i] = r
	}

	return ax
}
