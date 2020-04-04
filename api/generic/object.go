package generic

import (
	"github.com/milosgajdos/kraph/api"
)

// Object is generic API object
type Object struct {
	name  string
	kind  string
	ns    string
	uid   *UID
	links map[string]*Relation
}

// NewObject creates new Object and returns it
func NewObject(name, kind, ns string, uid *UID, links map[string]*Relation) *Object {
	return &Object{
		name:  name,
		kind:  kind,
		ns:    ns,
		uid:   uid,
		links: links,
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
func (o Object) UID() api.UID {
	return o.uid
}

// Links returns all links
func (o *Object) Links() []api.Link {
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

// Raw returns raw object
func (o Object) Raw() interface{} {
	return o
}
