package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Object is an API object
type Object struct {
	raw   unstructured.Unstructured
	links map[string]map[string]*ObjRef
}

// NewObject returns new kubernetes API object
func NewObject(raw unstructured.Unstructured) *Object {
	links := make(map[string]map[string]*ObjRef)

	for _, ref := range raw.GetOwnerReferences() {
		objRef := &ObjRef{
			name: ref.Name,
			kind: ref.Kind,
			uid:  string(ref.UID),
		}
		key := objRef.Name() + "/" + objRef.Kind()
		if links[key]["isOwned"] == nil {
			links[key] = make(map[string]*ObjRef)
		}
		links[key]["isOwned"] = objRef
	}

	return &Object{
		raw:   raw,
		links: links,
	}
}

// Name returns resource nam
func (o Object) Name() string {
	return strings.ToLower(o.raw.GetName())
}

// Kind returns object kind
func (o Object) Kind() string {
	return strings.ToLower(o.raw.GetKind())
}

// Namespace returns the namespace
func (o Object) Namespace() string {
	return strings.ToLower(o.raw.GetNamespace())
}

// UID returns object UID
func (o Object) UID() string {
	kind := o.Kind()
	uid := o.raw.GetUID()
	if len(uid) == 0 {
		return kind + "-" + o.Name()
	}

	return string(uid)
}

// Link links object to the given ref assigning the link the given relation
func (o *Object) Link(ref api.ObjRef, rel api.Relation) error {
	if o.links == nil {
		o.links = make(map[string]map[string]*ObjRef)
	}

	objRef := &ObjRef{
		name: ref.Name(),
		kind: ref.Kind(),
		uid:  ref.UID(),
	}

	key := objRef.Name() + "/" + objRef.Kind()
	if o.links[key][rel.String()] == nil {
		o.links[key] = make(map[string]*ObjRef)
	}

	if _, ok := o.links[key][rel.String()]; !ok {
		o.links[key][rel.String()] = objRef
	}

	return nil
}

// Links returns all object links
func (o Object) Links() []api.Link {
	var links []api.Link

	for _, rels := range o.links {
		for rel, ref := range rels {
			link := &Link{
				ref: ref,
				rel: &Relation{
					r: rel,
				},
			}
			links = append(links, link)
		}
	}

	return links
}

// Raw returns the raw Kubernetes API object
// The underlying type  of the returned object is unstructured.Unstructured
// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured#Unstructured
func (o Object) Raw() interface{} {
	return o.raw
}
