package generic

import (
	"github.com/milosgajdos/kraph/api"
)

// Object is generic API object
type Object struct {
	name  string
	kind  string
	ns    string
	uid   string
	links map[string]map[string]*ObjRef
}

// NewObject creates new Object and returns it
func NewObject(name, kind, ns, uid string) *Object {
	return &Object{
		name:  name,
		kind:  kind,
		ns:    ns,
		uid:   uid,
		links: make(map[string]map[string]*ObjRef),
	}
}

// Name returns object name
func (o Object) Name() string {
	return o.name
}

// Kind returns object kind
func (o Object) Kind() string {
	return o.kind
}

// Namespace returns object namespace
func (o Object) Namespace() string {
	return o.ns
}

// UID returns object uid
func (o Object) UID() string {
	return o.uid
}

// Link links the object to ref assigning the link the given relation
func (o *Object) Link(ref api.ObjRef, rel api.Relation) error {
	if o.links == nil {
		o.links = make(map[string]map[string]*ObjRef)
	}

	objRef := &ObjRef{
		name: ref.Name(),
		kind: ref.Kind(),
		uid:  ref.UID(),
	}

	key := objRef.name + "/" + objRef.kind
	if o.links[key][rel.String()] == nil {
		o.links[key] = make(map[string]*ObjRef)
	}

	if _, ok := o.links[key][rel.String()]; !ok {
		o.links[key][rel.String()] = objRef
	}

	return nil
}

// Links returns all links
func (m *Object) Links() []api.Link {
	var links []api.Link

	for _, rels := range m.links {
		for rel, obj := range rels {
			link := &Link{
				objRef: obj,
				relation: &Relation{
					r: rel,
				},
			}
			links = append(links, link)
		}
	}

	return links
}

// Raw returns raw object
func (o Object) Raw() interface{} {
	return o
}
