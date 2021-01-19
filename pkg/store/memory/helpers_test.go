package memory

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	objPath        = "testdata/objects.yaml"
	nodeResName    = "nodeResName"
	nodeResGroup   = "nodeResGroup"
	nodeResVersion = "nodeResVersion"
	nodeResKind    = "nodeResKind"
	nodeGID        = 123
	nodeID         = "testID"
	nodeName       = "testName"
	nodeNs         = "testNs"
)

func newTestResource(name, group, version, kind string, namespaced bool, opts api.Options) api.Resource {
	return generic.NewResource(name, group, version, kind, namespaced, opts)
}

func newTestObject(uid, name, ns string, res api.Resource, opts api.Options) api.Object {
	return generic.NewObject(uuid.NewFromString(uid), name, ns, res, opts)
}
