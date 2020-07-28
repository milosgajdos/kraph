package k8s

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/query"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func newTestAPI() *API {
	api := newAPI()

	for name, meta := range mock.Resources {
		groups := mock.ResourceData[name]["groups"]
		versions := mock.ResourceData[name]["versions"]
		for _, group := range groups {
			for _, version := range versions {
				var ns bool
				if len(meta["ns"]) > 0 {
					ns = true
				}
				res := Resource{
					ar: metav1.APIResource{
						Name:       name,
						Kind:       meta["kind"],
						Namespaced: ns,
					},
					gv: schema.GroupVersion{
						Group:   group,
						Version: version,
					},
				}
				api.AddResource(res)
				for _, path := range res.Paths() {
					api.AddResourceToPath(res, path)
				}
			}
		}
	}

	return api
}

func TestSource(t *testing.T) {
	api := newTestAPI()

	source := api.Source()

	expected := "kubernetes"
	if !strings.EqualFold(expected, source.String()) {
		t.Errorf("expected: %s, got: %s", expected, source.String())
	}
}

func TestResources(t *testing.T) {
	api := newTestAPI()

	resources := api.Resources()
	if len(resources) == 0 {
		t.Errorf("no resources found")
	}
}

func TestGetSimple(t *testing.T) {
	api := newTestAPI()

	for name := range mock.Resources {
		resources, err := api.Get(query.Name(name))
		if err != nil {
			t.Errorf("error querying name %s: %v", name, err)
		}

		for _, res := range resources {
			if res.Name() != name {
				t.Errorf("expected to get: %s, got: %s", name, res.Name())
			}
		}
		for _, group := range mock.ResourceData[name]["groups"] {
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
	api := newTestAPI()

	for name := range mock.Resources {
		for _, group := range mock.ResourceData[name]["groups"] {
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
	api := newTestAPI()

	for name := range mock.Resources {
		groups := mock.ResourceData[name]["groups"]
		versions := mock.ResourceData[name]["versions"]
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
