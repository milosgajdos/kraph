package k8s

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	apiCount = 7
	evenapi  = "even"
	oddapi   = "odd"
)

func makeTestAPI() API {
	api := API{
		resources:   make([]Resource, 0),
		resourceMap: make(map[string][]Resource),
	}

	for i := 0; i < apiCount; i++ {
		r := Resource{
			ar: metav1.APIResource{},
			gv: schema.GroupVersion{},
		}

		if i%2 == 0 {
			r.ar.Name = "even"
			api.resourceMap[evenapi] = append(api.resourceMap[evenapi], r)
		} else {
			r.ar.Name = "odd"
			api.resourceMap[oddapi] = append(api.resourceMap[oddapi], r)
		}

		api.resources = append(api.resources, r)
	}

	return api
}

func TestResources(t *testing.T) {
	api := makeTestAPI()

	resources := api.Resources("")

	if len(resources) != len(api.resources) {
		t.Errorf("expected %d API resources, got: %d", len(api.resources), len(resources))
	}

	oddResources := api.Resources(oddapi)
	if len(oddResources) != len(api.resourceMap[oddapi]) {
		t.Errorf("expected %d odd API resources, got: %d", len(api.resourceMap[oddapi]), len(oddResources))
	}

	evenResources := api.Resources(evenapi)
	if len(evenResources) != len(api.resourceMap[evenapi]) {
		t.Errorf("expected %d aven API resources, got: %d", len(api.resourceMap[evenapi]), len(oddResources))
	}
}
