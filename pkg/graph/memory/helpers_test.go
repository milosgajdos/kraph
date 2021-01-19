package memory

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
	"github.com/milosgajdos/kraph/pkg/api/types"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/graph"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
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

func makeTestAPIObjects(path string) (map[string]api.Object, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var testObjects []types.Object
	if err := yaml.Unmarshal(data, &testObjects); err != nil {
		return nil, err
	}

	objects := make(map[string]api.Object)

	for _, o := range testObjects {
		res := generic.NewResource(
			o.Resource.Name,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Kind,
			o.Resource.Namespaced,
			api.Options{Metadata: metadata.NewFromMap(o.Resource.Metadata)},
		)

		obj := generic.NewObject(uuid.NewFromString(o.UID), o.Name, o.Namespace, res, api.Options{Metadata: metadata.NewFromMap(o.Metadata)})

		for _, l := range o.Links {
			obj.Link(uuid.NewFromString(l.To), api.LinkOptions{Metadata: metadata.NewFromMap(l.Metadata)})
		}

		objects[o.UID] = obj
	}

	return objects, nil
}

func makeTestGraph(path string) (*WUG, error) {
	g, err := NewWUG("test", graph.Options{})
	if err != nil {
		return nil, err
	}

	objects, err := makeTestAPIObjects(path)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		n, err := g.NewNode(object)
		if err != nil {
			return nil, err
		}

		if err := g.AddNode(n); err != nil {
			return nil, err
		}

		for _, link := range object.Links() {
			object2 := objects[link.To().String()]

			n2, err := g.NewNode(object2)
			if err != nil {
				return nil, err
			}

			if err := g.AddNode(n2); err != nil {
				return nil, err
			}

			attrs := attrs.New()
			if relation, ok := link.Metadata().Get("relation").(string); ok {
				attrs.Set("relation", relation)
			}

			if _, err = g.Link(n.UID(), n2.UID(), graph.LinkOptions{Attrs: attrs}); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}
