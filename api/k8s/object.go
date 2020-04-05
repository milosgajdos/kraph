package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Object is an API object
type Object struct {
	name  string
	kind  string
	ns    string
	uid   *UID
	raw   unstructured.Unstructured
	links map[string]*Relation
}

// NewObject returns new kubernetes API object
func NewObject(raw unstructured.Unstructured) *Object {
	links := make(map[string]*Relation)

	name := strings.ToLower(raw.GetName())
	kind := strings.ToLower(raw.GetKind())
	ns := strings.ToLower(raw.GetNamespace())

	rawUID := string(raw.GetUID())
	if len(rawUID) == 0 {
		rawUID = kind + "-" + name
	}
	uid := &UID{uid: rawUID}

	for _, ref := range raw.GetOwnerReferences() {
		links[string(ref.UID)] = &Relation{rel: "isOwned"}
	}

	return &Object{
		name:  name,
		kind:  kind,
		ns:    ns,
		uid:   uid,
		raw:   raw,
		links: links,
	}
}

// Name returns resource nam
func (o Object) Name() string {
	return o.name
}

// Kind returns object kind
func (o Object) Kind() string {
	return o.kind
}

// Namespace returns the namespace
func (o Object) Namespace() string {
	return o.ns
}

// UID returns object UID
func (o Object) UID() api.UID {
	return o.uid
}

// Links returns all object links
func (o Object) Links() []api.Link {
	var links []api.Link

	for uid, rel := range o.links {
		link := &Link{
			to:  &UID{uid: uid},
			rel: rel,
		}
		links = append(links, link)
	}

	return links
}

// Raw returns the raw Kubernetes API object
// The underlying type  of the returned object is unstructured.Unstructured
// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured#Unstructured
func (o Object) Raw() interface{} {
	return o.raw
}
