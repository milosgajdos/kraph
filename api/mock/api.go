package mock

import (
	"github.com/milosgajdos/kraph/api/generic"
)

var (
	Resources = map[string]map[string]string{
		"foo": {
			"kind": "fooKind",
			"ns":   "foo",
		},
		"bar": {
			"kind": "barKind",
			"ns":   "",
		},
		"rnd": {
			"kind": "rndKind",
			"ns":   "rnd",
		},
	}
	ResourceData = map[string]map[string][]string{
		"foo": {
			"groups":   []string{"fooGroup", "wooGroup"},
			"versions": []string{"v1", "v2"},
		},
		"bar": {
			"groups":   []string{"barGroup", "carGroup"},
			"versions": []string{"v2"},
		},
		"rnd": {
			"groups":   []string{"rndGroup", "sndGroup"},
			"versions": []string{"v2", "v5", "v6"},
		},
	}
)

// API provides mock API
type API struct {
	*generic.API
}

// NewAPI returns new mock API
func NewAPI() *API {
	api := &API{
		API: generic.NewAPI(),
	}

	for name, meta := range Resources {
		groups := ResourceData[name]["groups"]
		versions := ResourceData[name]["versions"]
		for _, group := range groups {
			for _, version := range versions {
				var ns bool
				if len(meta["ns"]) > 0 {
					ns = true
				}
				res := generic.NewResource(name, meta["kind"], group, version, ns)
				api.AddResource(*res)
				for _, path := range res.Paths() {
					api.AddResourceToPath(*res, path)
				}
			}
		}
	}

	return api
}
