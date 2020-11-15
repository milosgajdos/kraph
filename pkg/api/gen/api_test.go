package gen

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/pkg/query"
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

	groups := make([]string, len(api.Resources()))
	versions := make([]string, len(api.Resources()))
	kinds := make([]string, len(api.Resources()))
	names := make([]string, len(api.Resources()))

	for i, r := range api.Resources() {
		groups[i] = r.Group()
		versions[i] = r.Version()
		kinds[i] = r.Kind()
		names[i] = r.Name()
	}

	for _, group := range groups {
		q := query.Build().Group(group, query.StringEqFunc(group))

		resources, err := api.Get(q)
		if err != nil {
			t.Errorf("error querying group %s: %v", group, err)
		}

		for _, r := range resources {
			if r.Group() != group {
				t.Errorf("expected to get group: %s, got: %s", group, r.Group())
			}
		}

		for _, version := range versions {
			q = q.Version(version, query.StringEqFunc(version))

			resources, err := api.Get(q)
			if err != nil {
				t.Errorf("error querying group/vresion %s/%s: %v", group, version, err)
			}

			for _, res := range resources {
				if res.Version() != version || res.Group() != group {
					t.Errorf("expected to get: %s/%s, got: %s/%s", group, version, res.Group(), res.Version())
				}
			}

			for _, kind := range kinds {
				q = q.Kind(kind, query.StringEqFunc(kind))

				resources, err := api.Get(q)
				if err != nil {
					t.Errorf("error querying group/version/kind: %s/%s/%s: %v", group, version, kind, err)
				}

				for _, res := range resources {
					if res.Kind() != kind || res.Version() != version || res.Group() != group {
						t.Errorf("expected to get: %s/%s/%s, got: %s/%s/%s", group, version, kind,
							res.Group(), res.Version(), res.Kind())
					}
				}

				for _, name := range names {
					q = q.Name(name, query.StringEqFunc(name))

					resources, err := api.Get(q)
					if err != nil {
						t.Errorf("error querying group/version/kind/name: %s/%s/%s/%s: %v", group, version, kind, name, err)
					}

					for _, res := range resources {
						if res.Name() != name || res.Kind() != kind || res.Version() != version || res.Group() != group {
							t.Errorf("expected to get: %s/%s/%s/%s, got: %s/%s/%s/%s", group, version, kind, name,
								res.Group(), res.Version(), res.Kind(), res.Name())
						}
					}

				}
			}
		}
	}
}
