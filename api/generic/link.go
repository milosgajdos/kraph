package generic

import "github.com/milosgajdos/kraph/api"

// ObjRef is API object reference used for linking API objects
type ObjRef struct {
	name string
	kind string
	uid  string
}

// NewObjRef returns generic object reference
func NewObjRef(name, kind, uid string) *ObjRef {
	return &ObjRef{
		name: name,
		kind: kind,
		uid:  uid,
	}
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

// NewRelation returns generic link relation
func NewRelation(r string) *Relation {
	return &Relation{
		r: r,
	}
}

// String returns relation description
func (r *Relation) String() string {
	return r.r
}

// Link defines API object relation
type Link struct {
	objRef   *ObjRef
	relation *Relation
}

// NewLink returns generic link
func NewLink(ref *ObjRef, rel *Relation) *Link {
	return &Link{
		objRef:   ref,
		relation: rel,
	}
}

// ObjRef returns link object reference
func (l *Link) To() api.ObjRef {
	return l.objRef
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.relation
}
