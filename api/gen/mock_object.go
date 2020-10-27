package gen

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/uuid"
)

// NewMockObject creates new mock API object and returns it
func NewMockObject(uid, name, ns string, res api.Resource) api.Object {
	return NewObject(uuid.NewFromString(uid), name, ns, res)
}
