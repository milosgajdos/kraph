package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/gen"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// OwnRel is k8s api object relation
	OwnRel = "isOwned"
)

// Object is kubernetes API object
type Object struct {
	*gen.Object
}

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
	uid := gen.NewUID(rawUID)

	obj := &Object{
		Object: gen.NewObject(uid, name, ns, res),
	}

	for _, ref := range raw.GetOwnerReferences() {
		//fmt.Printf("Object %s/%s/%s/%s owned by %s\n", obj.Resource().Version(), obj.Namespace(), obj.Resource().Kind(), obj.Name(), string(ref.UID))
		obj.Link(gen.NewUID(string(ref.UID)), gen.NewRelation(OwnRel))
	}

	return obj
}
