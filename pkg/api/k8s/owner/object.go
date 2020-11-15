package owner

import (
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ownRel = "own"
)

// Object is kubernetes API object
type Object struct {
	*gen.Object
}

// TODO: Link() returns error, make this func return tuple
// NewObject returns new kubernetes API object
func NewObject(res api.Resource, raw unstructured.Unstructured) *Object {
	name := strings.ToLower(raw.GetName())
	kind := strings.ToLower(raw.GetKind())

	ns := api.NsGlobal
	if res.Namespaced() {
		ns = strings.ToLower(raw.GetNamespace())
	}

	rawUID := string(raw.GetUID())
	if len(rawUID) == 0 {
		rawUID = kind + "-" + name
	}
	uid := uuid.NewFromString(rawUID)

	// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta
	m := metadata.New()
	m.Set("created_at", raw.GetCreationTimestamp().Time)
	m.Set("cluster_name", raw.GetClusterName())
	m.Set("labels", raw.GetLabels())
	m.Set("annotations", raw.GetAnnotations())

	obj := &Object{
		Object: gen.NewObject(uid, name, ns, res, api.Options{Metadata: m}),
	}

	for _, ref := range raw.GetOwnerReferences() {
		// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#OwnerReference
		m := metadata.New()
		m.Set("relation", ownRel)
		m.Set("controller", ref.Controller)

		obj.Link(uuid.NewFromString(string(ref.UID)), api.LinkOptions{Metadata: m})
	}

	return obj
}
