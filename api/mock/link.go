package mock

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/generic"
)

type Relation struct {
	*generic.Relation
}

func NewRelation(r string) api.Relation {
	return &Relation{
		Relation: generic.NewRelation(r),
	}
}
