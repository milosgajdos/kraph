package generic

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
)

// Top is Kubernetes API topology
type Top struct {
	// objects indexes all objects by their UID
	objects map[string]api.Object
	// index is a "search index" (ns/kind/name)
	index map[string]map[string]map[string]api.Object
}

// newTopology creates a new empty topology and returns it
func newTopology() *Top {
	return &Top{
		objects: make(map[string]api.Object),
		index:   make(map[string]map[string]map[string]api.Object),
	}
}

func (t Top) getNamespaceKindObjects(ns, kind string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Name != api.NameAll {
		object, ok := t.index[ns][kind][q.Name]
		if !ok {
			return objects, nil
		}

		objects = append(objects, object)

		return objects, nil
	}

	for _, object := range t.index[ns][kind] {
		objects = append(objects, object)
	}

	return objects, nil
}

func (t Top) getNamespaceObjects(ns string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Kind != api.KindAll {
		return t.getNamespaceKindObjects(ns, q.Kind, q)
	}

	for kind, _ := range t.index[ns] {
		objs, err := t.getNamespaceKindObjects(ns, kind, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

func (t Top) getAllNamespaceObjects(q query.Options) ([]api.Object, error) {
	var objects []api.Object

	for ns, _ := range t.index {
		objs, err := t.getNamespaceObjects(ns, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

// Objects returns all api objects in the tpoology
func (t Top) Objects() []api.Object {
	objects := make([]api.Object, len(t.objects))

	i := 0

	for _, object := range t.objects {
		objects[i] = object
		i++
	}

	return objects
}

// Get queries the mapped API objects and returns them
func (t Top) Get(opts ...query.Option) ([]api.Object, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var objects []api.Object

	if len(query.UID) > 0 {
		if obj, ok := t.objects[query.UID]; ok {
			objects = append(objects, obj)
		}
		return objects, nil
	}

	if len(query.Namespace) == 0 {
		return t.getAllNamespaceObjects(query)
	}

	objs, err := t.getNamespaceObjects(query.Namespace, query)
	if err != nil {
		return nil, err
	}
	objects = append(objects, objs...)

	return objects, nil
}
