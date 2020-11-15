package owner

import (
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Resource is API resource
type Resource struct {
	*gen.Resource
	ar metav1.APIResource
	gv schema.GroupVersion
}

// NewResource creates a new API resource and returns it.
func NewResource(ar metav1.APIResource, gv schema.GroupVersion, opts api.Options) *Resource {
	return &Resource{
		ar:       ar,
		gv:       gv,
		Resource: gen.NewResource(ar.Name, ar.Kind, gv.Group, gv.Version, ar.Namespaced, opts),
	}
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
