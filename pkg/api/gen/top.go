package gen

import (
	"sync"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Top is generic API topology
type Top struct {
	// api is the source of topology
	api api.API
	// objects stores all objects by their UID
	objects map[string]api.Object
	// index is the topology "search index" (ns/kind/name)
	index map[string]map[string]map[string]api.Object
	// mu synchronizes access to Top
	mu *sync.RWMutex
}

// NewTop creates a new empty topology and returns it
func NewTop(a api.API) *Top {
	return &Top{
		api:     a,
		objects: make(map[string]api.Object),
		index:   make(map[string]map[string]map[string]api.Object),
		mu:      &sync.RWMutex{},
	}
}

// API returns topology API source
func (t Top) API() api.API {
	return t.api
}

// add adds a new object to topology
func (t *Top) add(o api.Object) error {
	t.objects[o.UID().String()] = o

	ns := o.Namespace()

	if t.index[ns] == nil {
		t.index[ns] = make(map[string]map[string]api.Object)
	}

	kind := o.Resource().Kind()

	if t.index[ns][kind] == nil {
		t.index[ns][kind] = make(map[string]api.Object)
	}

	name := o.Name()

	t.index[ns][kind][name] = o

	return nil
}

// Add adds o to the topology using on the provided options.
// If an object already exists in the topology and MergeLinks option is enabled
// the existing object links are merged with the links of o.
func (t *Top) Add(o api.Object, opts api.AddOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	obj, ok := t.objects[o.UID().String()]
	if !ok {
		if err := t.add(o); err != nil {
			return err
		}

		for _, l := range o.Links() {
			lopts := api.LinkOptions{
				UID:      l.UID(),
				Multi:    opts.MultiLink,
				Metadata: l.Metadata(),
			}

			// link: to -> o
			if to, ok := t.objects[l.To().String()]; ok {
				if err := to.Link(o.UID(), lopts); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if opts.MergeLinks {
		for _, l := range o.Links() {
			lopts := api.LinkOptions{
				UID:      l.UID(),
				Multi:    opts.MultiLink,
				Metadata: l.Metadata(),
			}

			// link: o -> to
			if err := obj.Link(l.To(), lopts); err != nil {
				return err
			}

			// link: to -> o
			if to, ok := t.objects[l.To().String()]; ok {
				if err := to.Link(obj.UID(), lopts); err != nil {
					return err
				}
			}
		}
	}

	return nil
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
	t.mu.RLock()
	defer t.mu.RUnlock()

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
	t.mu.RLock()
	defer t.mu.RUnlock()

	var objects []api.Object

	if m := q.Matcher().UID(); m != nil {
		// NOTE: when we fail to type-switch, we fall through to query.MatchAny
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
