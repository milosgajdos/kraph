package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Object is a generic API object
type Object struct {
	uid   uuid.UID
	name  string
	ns    string
	res   api.Resource
	links map[string]api.Link
}

// NewObject creates a new Object and returns it
func NewObject(uid uuid.UID, name, ns string, res api.Resource) *Object {
	return &Object{
		uid:   uid,
		res:   res,
		ns:    ns,
		name:  name,
		links: make(map[string]api.Link),
	}
}

// UID returns object uid
func (o Object) UID() uuid.UID {
	return o.uid
}

// Name returns object name
func (o Object) Name() string {
	return o.name
}

// Namespace returns object namespace
func (o Object) Namespace() string {
	return o.ns
}

// Resource returns API resource the object is an instance of
func (o Object) Resource() api.Resource {
	return o.res
}

// Link links the object to another object
func (o *Object) Link(to uuid.UID, rel api.Relation) {
	link := NewLink(o.uid, to, rel)

	o.links[link.UID().String()] = link
}

// Links returns a slice of all object links
func (o *Object) Links() []api.Link {
	var links []api.Link

	for _, link := range o.links {
		links = append(links, link)
	}

	return links
}
