package mock

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/generic"
)

type Object struct {
	*generic.Object
}

func NewObject(name, kind, ns, uid string, links map[string]api.Relation) api.Object {
	return &Object{
		Object: generic.NewObject(name, kind, ns, generic.NewUID(uid), links),
	}
}
