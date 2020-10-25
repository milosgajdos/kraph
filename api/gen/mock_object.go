package gen

import "github.com/milosgajdos/kraph/api"

// NewMockObject creates new mock API object and returns it
func NewMockObject(uid, name, ns string, res api.Resource) api.Object {
	return NewObject(NewUID(uid), name, ns, res)
}
