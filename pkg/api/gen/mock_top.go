package gen

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/types"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockTop returns mock Top from objects and resrouces
// from given filesystem paths and returns it.
func NewMockTop(a api.API, objPath string) (*Top, error) {
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		return nil, err
	}

	var objects []types.Object
	if err := yaml.Unmarshal(data, &objects); err != nil {
		return nil, err
	}

	top := NewTop(a)

	for _, o := range objects {
		r := &Resource{
			name:       o.Resource.Name,
			kind:       o.Resource.Kind,
			group:      o.Resource.Group,
			version:    o.Resource.Version,
			namespaced: o.Resource.Namespaced,
		}

		obj := &Object{
			uid:    uuid.NewFromString(o.UID),
			name:   o.Name,
			ns:     o.Namespace,
			res:    r,
			links:  make(map[string]api.Link),
			olinks: make(map[string]struct{}),
		}

		for _, l := range o.Links {
			m := metadata.New()

			for k, v := range l.Metadata {
				m.Set(k, v)
			}

			obj.links[l.UID] = Link{
				uid:  uuid.NewFromString(l.UID),
				from: uuid.NewFromString(l.From),
				to:   uuid.NewFromString(l.To),
				opts: api.LinkOptions{Metadata: m}}
			obj.olinks[l.To] = struct{}{}
		}

		if err := top.Add(obj, api.AddOptions{}); err != nil {
			return nil, err
		}
	}

	return top, nil
}
