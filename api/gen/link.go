package gen

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/uuid"
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

// Link links API object to another API object
type Link struct {
	uid  uuid.UID
	from uuid.UID
	to   uuid.UID
	rel  api.Relation
}

// NewLink returns a new link between API objects
func NewLink(from, to uuid.UID, rel api.Relation) *Link {
	return &Link{
		uid:  uuid.New(),
		from: from,
		to:   to,
		rel:  rel,
	}
}

// UID returns link uid
func (l *Link) UID() uuid.UID {
	return l.uid
}

// From returns linking object reference
func (l *Link) From() uuid.UID {
	return l.from
}

// To returns link object reference
func (l *Link) To() uuid.UID {
	return l.to
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.rel
}
