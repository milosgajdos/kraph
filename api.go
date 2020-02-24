package kraph

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Resource is API resource
type Resource struct {
	ar metav1.APIResource
	gv schema.GroupVersion
}

// Paths returns all possible variations of the resource paths
func (r Resource) Paths() []string {
	// WTF: SingularName is often empty string!
	singularName := r.ar.SingularName
	if singularName == "" {
		singularName = r.ar.Kind
	}
	resNames := []string{strings.ToLower(r.ar.Name), strings.ToLower(singularName)}
	resNames = append(resNames, r.ar.ShortNames...)

	var names []string
	for _, name := range resNames {
		names = append(names,
			name,
			strings.Join([]string{name, r.gv.Group}, "/"),
			strings.Join([]string{name, r.gv.Group, r.gv.Version}, "/"),
		)
	}

	return names
}

// API is API resource map
type API struct {
	resources   []Resource
	resourceMap map[string][]Resource
}

// Resources returns API resources
func (a *API) Resources() []Resource {
	resources := make([]Resource, len(a.resources))

	for i, r := range a.resources {
		resources[i] = r
	}

	return resources
}

// Lookup looks up all API resources for the given API name and returns them
func (a *API) Lookup(name string) []Resource {
	var resources []Resource

	if a.resourceMap == nil {
		a.resourceMap = make(map[string][]Resource)
		return resources
	}

	if apiResources, ok := a.resourceMap[name]; ok {
		resources = make([]Resource, len(apiResources))

		for i, r := range apiResources {
			resources[i] = r
		}
	}

	return resources
}
