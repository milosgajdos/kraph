package gen

import (
	"github.com/google/uuid"
	"github.com/milosgajdos/kraph/api"
)

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

// Link links API object to another API object
type Link struct {
	uid  api.UID
	from api.UID
	to   api.UID
	rel  api.Relation
}

// NewLink returns a new link between API objects
func NewLink(from, to api.UID, rel api.Relation) *Link {
	return &Link{
		uid:  &UID{uid: uuid.New().String()},
		from: from,
		to:   to,
		rel:  rel,
	}
}

// UID returns link uid
func (l *Link) UID() api.UID {
	return l.uid
}

// From returns linking object reference
func (l *Link) From() api.UID {
	return l.from
}

// To returns link object reference
func (l *Link) To() api.UID {
	return l.to
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.rel
}
