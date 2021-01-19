package generic

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/types"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockTop returns mock Top from objects and resrouces
// from given path and returns it.
func NewMockTop(a api.API, path string) (*Top, error) {
	data, err := ioutil.ReadFile(path)
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
			olinks: make(map[string]api.Link),
			md:     metadata.New(),
		}

		for _, l := range o.Links {
			m := metadata.New()

			for k, v := range l.Metadata {
				m.Set(k, v)
			}

			link := Link{
				uid:  uuid.NewFromString(l.UID),
				from: uuid.NewFromString(l.From),
				to:   uuid.NewFromString(l.To),
				md:   metadata.NewCopyFrom(m),
			}

			obj.links[l.UID] = link
			obj.olinks[l.To] = link
		}

		if err := top.Add(obj, api.AddOptions{}); err != nil {
			return nil, err
		}
	}

	return top, nil
}
