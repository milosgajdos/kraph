package mock

import (
	"testing"

	"github.com/milosgajdos/kraph/query"
)

func TestResources(t *testing.T) {
	api := NewAPI()

	resources := api.Resources()
	if len(resources) == 0 {
		t.Errorf("no resources found")
	}
}

func TestGetSimple(t *testing.T) {
	api := NewAPI()

	for name, _ := range Resources {
		resources, err := api.Get(query.Name(name))
		if err != nil {
			t.Errorf("error querying name %s: %v", name, err)
		}

		for _, res := range resources {
			if res.Name() != name {
				t.Errorf("expected to get: %s, got: %s", name, res.Name())
			}
		}
		for _, group := range ResourceData[name]["groups"] {
			resources, err := api.Get(query.Group(group))
			if err != nil {
				t.Errorf("error querying name %s: %v", name, err)
			}

			for _, res := range resources {
				if res.Group() != group {
					t.Errorf("expected to get: %s, got: %s", group, res.Group())
				}
			}
		}
	}
}

func TestGetNameGroup(t *testing.T) {
	api := NewAPI()

	for name, _ := range Resources {
		for _, group := range ResourceData[name]["groups"] {
			resources, err := api.Get(query.Name(name), query.Group(group))
			if err != nil {
				t.Errorf("error querying name/group: %s/%s: %v", name, group, err)
			}

			for _, res := range resources {
				if res.Name() != name || res.Group() != group {
					t.Errorf("expected to get: %s/%s, got: %s/%s", group, name, res.Name(), res.Group())
				}
			}
		}
	}
}

func TestGetNameGroupVersion(t *testing.T) {
	api := NewAPI()

	for name, _ := range Resources {
		groups := ResourceData[name]["groups"]
		versions := ResourceData[name]["versions"]
		for _, group := range groups {
			for _, version := range versions {
				resources, err := api.Get(query.Name(name), query.Group(group), query.Version(version))
				if err != nil {
					t.Errorf("error querying name/group/version: %s/%s/%s: %v", name, group, version, err)
				}

				for _, res := range resources {
					if res.Name() != name || res.Group() != group || res.Version() != version {
						t.Errorf("expected to get: %s/%s/%s, got: %s/%s/%s", group, name, version,
							res.Name(), res.Group(), res.Version())
					}
				}
			}
		}
	}
}
