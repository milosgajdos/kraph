package gen

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/uuid"
)

// Top is generic API topology
type Top struct {
	// objects indexes all objects by their UID
	objects map[string]api.Object
	// index is a "search index" (ns/kind/name)
	index map[string]map[string]map[string]api.Object
}

// NewTop creates a new empty topology and returns it
func NewTop() *Top {
	return &Top{
		objects: make(map[string]api.Object),
		index:   make(map[string]map[string]map[string]api.Object),
	}
}

// Add adds an Object to the topology
func (t *Top) Add(o api.Object) {
	if _, ok := t.objects[o.UID().String()]; !ok {
		t.objects[o.UID().String()] = o

		if t.index[o.Namespace()] == nil {
			t.index[o.Namespace()] = make(map[string]map[string]api.Object)
		}

		kind := o.Resource().Kind()

		if t.index[o.Namespace()][kind] == nil {
			t.index[o.Namespace()][kind] = make(map[string]api.Object)
		}

		t.index[o.Namespace()][kind][o.Name()] = o
	}
}

func (t Top) getNamespaceKindObjects(ns, kind string, q *query.Query) ([]api.Object, error) {
	var objects []api.Object

	if m := q.Matcher().Name(); m != nil {
		switch name := m.Value().(type) {
		case string:
			if len(name) > 0 {
				object, ok := t.index[ns][kind][name]
				if !ok {
					return objects, nil
				}

				objects = append(objects, object)
			}
		case query.MatchVal:
			if name == query.MatchAny {
				for _, object := range t.index[ns][kind] {
					objects = append(objects, object)
				}
			}
		}
	}

	return objects, nil
}

func (t Top) getNamespaceObjects(ns string, q *query.Query) ([]api.Object, error) {
	var objects []api.Object

	if m := q.Matcher().Kind(); m != nil {
		switch kind := m.Value().(type) {
		case string:
			if len(kind) > 0 {
				return t.getNamespaceKindObjects(ns, kind, q)
			}
		case query.MatchVal:
			if kind == query.MatchAny {
				for kind := range t.index[ns] {
					objs, err := t.getNamespaceKindObjects(ns, kind, q)
					if err != nil {
						return nil, err
					}
					objects = append(objects, objs...)
				}
			}
		}
	}

	return objects, nil
}

func (t Top) getAllNamespacedObjects(q *query.Query) ([]api.Object, error) {
	var objects []api.Object

	for ns := range t.index {
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

// Get queries the mapped API objects and returns the results
func (t Top) Get(q *query.Query) ([]api.Object, error) {
	var objects []api.Object

	if m := q.Matcher().UID(); m != nil {
		// TODO: when we fail to type-switch, we fall through to MatchAny
		// Should this logic change? i.e. we can see that UID matcher is NOT nil
		// but the provided value does not type-switch; maybe we should return
		// some ErrInvalidUID error or something; this would avoid a lot of missteps
		if uid, ok := m.Value().(uuid.UID); ok && len(uid.String()) > 0 {
			if obj, ok := t.objects[uid.String()]; ok {
				objects = append(objects, obj)
			}
			return objects, nil
		}
	}

	var ns string

	if m := q.Matcher().Namespace(); m != nil {
		switch v := m.Value().(type) {
		case string:
			ns = v
		case query.MatchVal:
			return t.getAllNamespacedObjects(q)
		}
	}

	if len(ns) > 0 {
		objs, err := t.getNamespaceObjects(ns, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}
