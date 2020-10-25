package gen

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/query"
)

const (
	resPath = "seeds/resources.yaml"
)

func TestSource(t *testing.T) {
	api, err := NewMockAPI(resPath)
	if err != nil {
		t.Errorf("failed to create mock API: %v", err)
		return
	}

	if src := api.Source(); !strings.EqualFold(resPath, src.String()) {
		t.Errorf("expected: %s, got: %s", resPath, src.String())
	}
}

func TestResources(t *testing.T) {
	api, err := NewMockAPI(resPath)
	if err != nil {
		t.Errorf("failed to create mock API: %v", err)
		return
	}

	resources := api.Resources()
	if len(resources) == 0 {
		t.Errorf("no resources found")
	}
}

func TestAPIGet(t *testing.T) {
	api, err := NewMockAPI(resPath)
	if err != nil {
		t.Errorf("failed to create mock API: %v", err)
		return
	}

	names := make([]string, len(api.Resources()))
	groups := make([]string, len(api.Resources()))
	versions := make([]string, len(api.Resources()))
	for i, r := range api.Resources() {
		names[i] = r.Name()
		groups[i] = r.Group()
		versions[i] = r.Version()
	}

	for _, name := range names {
		q := query.Build().Name(name, query.StringEqFunc(name))

		resources, err := api.Get(q)
		if err != nil {
			t.Errorf("error querying name %s: %v", name, err)
		}

		for _, r := range resources {
			if r.Name() != name {
				t.Errorf("expected to get: %s, got: %s", name, r.Name())
			}
		}

		for _, group := range groups {
			q = q.Group(group, query.StringEqFunc(group))

			resources, err := api.Get(q)
			if err != nil {
				t.Errorf("error querying group/name %s/%s: %v", group, name, err)
			}

			for _, res := range resources {
				if res.Name() != name || res.Group() != group {
					t.Errorf("expected to get: %s/%s, got: %s/%s", group, name, res.Group(), res.Name())
				}
			}

			for _, version := range versions {
				q = q.Version(version, query.StringEqFunc(version))

				resources, err := api.Get(q)
				if err != nil {
					t.Errorf("error querying group/version/name: %s/%s/%s: %v", group, version, name, err)
				}

				for _, res := range resources {
					if res.Name() != name || res.Group() != group || res.Version() != version {
						t.Errorf("expected to get: %s/%s/%s, got: %s/%s/%s", group, version, name,
							res.Group(), res.Version(), res.Name())
					}
				}
			}
		}
	}
}
