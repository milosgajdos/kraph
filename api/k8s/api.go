package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// Object is API object i.e. instance of API resource
type Object struct {
	obj unstructured.Unstructured
}

// Name returns resource nam
func (o Object) Name() string {
	return strings.ToLower(o.obj.GetKind()) + "-" + strings.ToLower(o.obj.GetName())
}

// Kind returns object kind
func (o Object) Kind() string {
	return strings.ToLower(o.obj.GetKind())
}

// Namespace returns the namespace
func (o Object) Namespace() string {
	return o.obj.GetNamespace()
}

// UID returns object UID
func (o Object) UID() types.UID {
	kind := strings.ToLower(o.obj.GetKind())
	uid := o.obj.GetUID()
	if uid == "" {
		uid = types.UID(kind + "-" + strings.ToLower(o.obj.GetName()))
	}

	return uid
}

// Raw returns the raw API bjoect
func (o *Object) Raw() interface{} {
	return o.obj
}

// Resource is API resource
type Resource struct {
	ar metav1.APIResource
	gv schema.GroupVersion
}

// Name returns the name of the resource
func (r Resource) Name() string {
	return r.ar.Name
}

// Kind returns resource kind
func (r Resource) Kind() string {
	return r.ar.Kind
}

// Group returns the API group of the resource
func (r Resource) Group() string {
	return r.ar.Group
}

// Version returns the version of the resource
func (r Resource) Version() string {
	return r.gv.Version
}

// Namespaced returns true if the resource is namespaced
func (r Resource) Namespaced() bool {
	return r.ar.Namespaced
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
func (a *API) Resources() []api.Resource {
	resources := make([]api.Resource, len(a.resources))

	for i, r := range a.resources {
		resources[i] = r
	}

	return resources
}

// Lookup looks up all API resources for the given API name and returns them
func (a *API) Lookup(name string) []api.Resource {
	var resources []api.Resource

	if a.resourceMap == nil {
		a.resourceMap = make(map[string][]Resource)
		return resources
	}

	if apiResources, ok := a.resourceMap[name]; ok {
		resources = make([]api.Resource, len(apiResources))

		for i, r := range apiResources {
			resources[i] = r
		}
	}

	return resources
}
