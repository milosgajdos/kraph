package mock

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/generic"
)

func NewObject(name, kind, ns, uid string, links map[string]api.Relation) api.Object {
	return generic.NewObject(name, kind, ns, generic.NewUID(uid), links)
}
