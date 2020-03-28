package k8s

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	corev1 "k8s.io/api/core/v1"
)

// Top is Kubernetes API topology
type Top map[string]map[string]map[string]api.Object

func (t Top) getNamespaceKindObjects(ns, kind string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Name != NameAll {
		object, ok := t[ns][kind][q.Name]
		if !ok {
			return objects, nil
		}
		objects = append(objects, object)
		return objects, nil
	}

	for _, obj := range t[ns][kind] {
		objects = append(objects, obj)
	}

	return objects, nil
}

func (t Top) getNamespaceObjects(ns string, q query.Options) ([]api.Object, error) {
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

func (t Top) getAllNamespaceObjects(q query.Options) ([]api.Object, error) {
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
func (t Top) Get(opts ...query.Option) ([]api.Object, error) {
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
