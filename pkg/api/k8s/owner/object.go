package owner

import (
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	relation = "own"
)

// Object is kubernetes API object
type Object struct {
	*generic.Object
}

// NewObject returns new kubernetes API object
func NewObject(res api.Resource, raw unstructured.Unstructured) (*Object, error) {
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

	if !raw.GetCreationTimestamp().Time.IsZero() {
		m.Set("created_at", raw.GetCreationTimestamp().Time.String())
	}

	if raw.GetClusterName() != "" {
		m.Set("cluster_name", raw.GetClusterName())
	}

	for k, v := range raw.GetLabels() {
		if v != "" {
			m.Set(k, v)
		}
	}

	for k, v := range raw.GetAnnotations() {
		if v != "" {
			m.Set(k, v)
		}
	}

	obj := &Object{
		Object: generic.NewObject(uid, name, ns, res, api.Options{Metadata: m}),
	}

	for _, ref := range raw.GetOwnerReferences() {
		// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#OwnerReference
		m := metadata.New()
		m.Set("relation", relation)

		if err := obj.Link(uuid.NewFromString(string(ref.UID)), api.LinkOptions{Metadata: m}); err != nil {
			return nil, err
		}
	}

	return obj, nil
}
