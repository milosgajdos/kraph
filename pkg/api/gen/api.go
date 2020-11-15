package gen

import (
	"sync"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/query"
)

// API is a generic API
type API struct {
	// source is API source
	source api.Source
	// resources stores discovered API resources
	// indexed as grou/version/kind
	resources map[string]map[string]map[string]api.Resource
	// mu synchronizes access to API
	mu *sync.RWMutex
}

// NewAPI returns new K8s API object
func NewAPI(src api.Source) *API {
	return &API{
		source:    src,
		resources: make(map[string]map[string]map[string]api.Resource),
		mu:        &sync.RWMutex{},
	}
}

// Add adds resource to the API
func (a *API) Add(r api.Resource, opts api.AddOptions) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	group := r.Group()

	if a.resources[group] == nil {
		a.resources[group] = make(map[string]map[string]api.Resource)
	}

	version := r.Version()

	if a.resources[group][version] == nil {
		a.resources[group][version] = make(map[string]api.Resource)
	}

	kind := r.Kind()

	a.resources[group][version][kind] = r

	return nil
}

// Source returns API source
func (a API) Source() api.Source {
	return a.source
}

// Resources returns all API resources
func (a API) Resources() []api.Resource {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var resources []api.Resource

	for _, groups := range a.resources {
		for _, versions := range groups {
			for _, r := range versions {
				resources = append(resources, r)
			}
		}
	}

	return resources
}

func matchName(r api.Resource, q *query.Query) bool {
	if m := q.Matcher().Name(); m != nil {
		switch name := m.Value().(type) {
		case string:
			if len(name) > 0 {
				return name == r.Name()
			}
		case query.MatchVal:
			if name == query.MatchAny {
				return true
			}
		}
	}

	return true
}

func (a API) getGroupVersionResources(group, version string, q *query.Query) ([]api.Resource, error) {
	var resources []api.Resource

	if m := q.Matcher().Kind(); m != nil {
		switch kind := m.Value().(type) {
		case string:
			if len(kind) > 0 {
				r, ok := a.resources[group][version][kind]
				if !ok {
					return resources, nil
				}

				if matchName(r, q) {
					resources = append(resources, r)
				}
			}
		case query.MatchVal:
			if kind == query.MatchAny {
				for kind := range a.resources[group][version] {
					r := a.resources[group][version][kind]
					if matchName(r, q) {
						resources = append(resources, r)
					}
				}
			}
		}
	}

	return resources, nil
}

func (a API) getGroupResources(group string, q *query.Query) ([]api.Resource, error) {
	var resources []api.Resource

	if m := q.Matcher().Version(); m != nil {
		switch version := m.Value().(type) {
		case string:
			if len(version) > 0 {
				return a.getGroupVersionResources(group, version, q)
			}
		case query.MatchVal:
			if version == query.MatchAny {
				for version := range a.resources[group] {
					rx, err := a.getGroupVersionResources(group, version, q)
					if err != nil {
						return nil, err
					}
					resources = append(resources, rx...)
				}
			}
		}
	}

	return resources, nil
}

func (a API) getAllGroupedResources(q *query.Query) ([]api.Resource, error) {
	var resources []api.Resource

	for g := range a.resources {
		rx, err := a.getGroupResources(g, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}

// Get returns all API resources matching the given query
func (a API) Get(q *query.Query) ([]api.Resource, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var resources []api.Resource

	var group string

	if m := q.Matcher().Group(); m != nil {
		switch v := m.Value().(type) {
		case string:
			group = v
		case query.MatchVal:
			return a.getAllGroupedResources(q)
		}
	}

	if len(group) > 0 {
		rx, err := a.getGroupResources(group, q)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rx...)
	}

	return resources, nil
}
