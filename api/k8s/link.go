package k8s

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/gen"
)

// Relation is link relation
type Relation struct {
	*gen.Relation
}

// NewRelation returns gen new relation
func NewRelation(r string) *Relation {
	return &Relation{
		Relation: gen.NewRelation(r),
	}
}

// UID implements API object UID
type UID struct {
	*gen.UID
}

// NewUID returns new UID
func NewUID(uid string) *UID {
	return &UID{
		UID: gen.NewUID(uid),
	}
}

// Link defines API object relation
type Link struct {
	*gen.Link
}

// NewLink returns new link
func NewLink(from, to api.UID, rel api.Relation) *Link {
	return &Link{
		Link: gen.NewLink(from, to, rel),
	}
}
