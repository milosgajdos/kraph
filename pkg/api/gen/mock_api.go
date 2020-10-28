package gen

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/types"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockAPI returns mock API from given path and returns it
func NewMockAPI(resPath string) (api.API, error) {
	data, err := ioutil.ReadFile(resPath)
	if err != nil {
		return nil, err
	}

	var resources []types.Resource
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	api := NewAPI(resPath)

	for _, r := range resources {
		m := &Resource{
			name:       r.Name,
			kind:       r.Kind,
			group:      r.Group,
			version:    r.Version,
			namespaced: r.Namespaced,
		}
		api.AddResource(m)
		for _, path := range m.Paths() {
			api.IndexPath(m, path)
		}
	}

	return api, nil
}

// NewMockTop returns mock Top from objects and resrouces
// from given filesystem paths and returns it.
func NewMockTop(objPath string) (api.Top, error) {
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		return nil, err
	}

	var objects []types.Object
	if err := yaml.Unmarshal(data, &objects); err != nil {
		return nil, err
	}

	top := NewTop()

	for _, o := range objects {
		r := &Resource{
			name:       o.Resource.Name,
			kind:       o.Resource.Kind,
			group:      o.Resource.Group,
			version:    o.Resource.Version,
			namespaced: o.Resource.Namespaced,
		}

		links := make(map[string]api.Link)
		for _, l := range o.Links {
			links[l.UID] = &Link{
				uid:  uuid.NewFromString(l.UID),
				from: uuid.NewFromString(l.From),
				to:   uuid.NewFromString(l.To),
				rel:  NewRelation(l.Relation),
			}
		}

		m := &Object{
			uid:   uuid.NewFromString(o.UID),
			name:  o.Name,
			ns:    o.Namespace,
			res:   r,
			links: links,
		}
		top.Add(m)
	}

	return top, nil
}
