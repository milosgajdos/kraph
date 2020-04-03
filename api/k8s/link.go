package k8s

import (
	"github.com/milosgajdos/kraph/api"
)

// Relation is link relation
type Relation struct {
	r string
}

// String returns relation description
func (r *Relation) String() string {
	return r.r
}

// UID implements API object UID
type UID struct {
	uid string
}

// String returns API Object UID as string
func (u *UID) String() string {
	return u.uid
}

// Link defines API object relation
type Link struct {
	to  *UID
	rel *Relation
}

// To returns linked object reference
func (l *Link) To() api.UID {
	return l.to
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.rel
}
