package gen

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/types"
)

// NewMockAPI returns mock API from given path and returns it
func NewMockAPI(resPath string) (*API, error) {
	data, err := ioutil.ReadFile(resPath)
	if err != nil {
		return nil, err
	}

	var resources []types.Resource
	if err := yaml.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	s := NewSource(resPath)
	a := NewAPI(s)

	for _, r := range resources {
		m := &Resource{
			name:       r.Name,
			kind:       r.Kind,
			group:      r.Group,
			version:    r.Version,
			namespaced: r.Namespaced,
		}

		if err := a.Add(m, api.AddOptions{}); err != nil {
			return nil, err
		}
	}

	return a, nil
}
