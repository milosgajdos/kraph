package k8s

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
	return r.gv.Group
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
	// WTF: SingularName is often an empty string!
	// TODO: figure this out; but for now let's set it to Kind
	singularName := r.ar.SingularName
	if len(singularName) == 0 {
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
