package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockObject creates new mock API object and returns it
func NewMockObject(uid, name, ns string, res api.Resource, opts api.Options) api.Object {
	return NewObject(uuid.NewFromString(uid), name, ns, res, opts)
}
