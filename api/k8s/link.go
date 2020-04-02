package k8s

import (
	"github.com/milosgajdos/kraph/api"
)

// ObjRef is object reference used for linking API objects
type ObjRef struct {
	name string
	kind string
	uid  string
}

// Name of the API object reference
func (r ObjRef) Name() string {
	return r.name
}

// Kind of the API object references
func (r ObjRef) Kind() string {
	return r.kind
}

// UID of the API object reference
func (r ObjRef) UID() string {
	return r.uid
}

// Relation is link relation
type Relation struct {
	r string
}

// String returns relation description
func (r *Relation) String() string {
	return r.r
}

// Link defines API object relation
type Link struct {
	ref *ObjRef
	rel *Relation
}

// To returns linked object reference
func (l *Link) To() api.ObjRef {
	return l.ref
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.rel
}
