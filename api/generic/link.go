package generic

import "github.com/milosgajdos/kraph/api"

// UID implements API object UID
type UID struct {
	uid string
}

// NewUID returns new UID
func NewUID(uid string) *UID {
	return &UID{
		uid: uid,
	}
}

// String returns API Object UID as string
func (u *UID) String() string {
	return u.uid
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

// Link links API object to another API object
type Link struct {
	to  *UID
	rel api.Relation
}

// NewLink returns generic link
func NewLink(to *UID, rel *Relation) *Link {
	return &Link{
		to:  to,
		rel: rel,
	}
}

// Ref returns link object reference
func (l *Link) To() api.UID {
	return l.to
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.rel
}
