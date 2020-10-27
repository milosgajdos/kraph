package k8s

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/uuid"
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
	uuid.UID
}

// NewUID returns new UID
func NewUID(uid string) *UID {
	return &UID{
		UID: uuid.NewFromString(uid),
	}
}

// Link defines API object relation
type Link struct {
	*gen.Link
}

// NewLink returns new link
func NewLink(from, to uuid.UID, rel api.Relation) *Link {
	return &Link{
		Link: gen.NewLink(from, to, rel),
	}
}
