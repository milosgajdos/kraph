package k8s

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	corev1 "k8s.io/api/core/v1"
)

// top is api topology
type top map[string]map[string]map[string]*Object

func (t top) getNamespaceKindObjects(ns, kind string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Name != NameAll {
		objects = append(objects, t[ns][kind][q.Name])
		return objects, nil
	}

	for name, _ := range t[ns][kind] {
		objects = append(objects, t[ns][kind][name])
	}

	return objects, nil
}

func (t top) getNamespaceObjects(ns string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Kind != KindAll {
		return t.getNamespaceKindObjects(ns, q.Kind, q)
	}

	for kind, _ := range t[ns] {
		objs, err := t.getNamespaceKindObjects(ns, kind, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

func (t top) getAllNamespaceObjects(q query.Options) ([]api.Object, error) {
	var objects []api.Object

	for ns, _ := range t {
		objs, err := t.getNamespaceObjects(ns, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

// Get queries the mapped API objects and returns them
// It never returns error
func (t top) Get(opts ...query.Option) ([]api.Object, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var objects []api.Object

	if query.Namespace == corev1.NamespaceAll {
		return t.getAllNamespaceObjects(query)
	}

	objs, err := t.getNamespaceObjects(query.Namespace, query)
	if err != nil {
		return nil, err
	}
	objects = append(objects, objs...)

	return objects, nil
}
