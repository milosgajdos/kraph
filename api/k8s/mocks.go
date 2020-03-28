package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/milosgajdos/kraph/api"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	testdiscclient "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	testdynclient "k8s.io/client-go/dynamic/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
)

var (
	MockAPIResCount = 9
	MockAPIGroups   = []string{"even", "odd"}
)

func MockAPI() API {
	api := API{
		resources:   make([]Resource, 0),
		resourceMap: make(map[string][]Resource),
	}

	for i := 0; i < MockAPIResCount; i++ {
		r := Resource{
			ar: metav1.APIResource{},
			gv: schema.GroupVersion{},
		}

		if i%2 == 0 {
			r.ar.Name, r.ar.SingularName = "evenRes", "er"
			r.gv.Group = MockAPIGroups[0]
		} else {
			r.ar.Name, r.ar.SingularName = "oddRes", "or"
			r.gv.Group = MockAPIGroups[1]
		}

		r.gv.Version = fmt.Sprintf("v%d", i)

		api.resources = append(api.resources, r)
		for _, path := range r.Paths() {
			api.resourceMap[path] = append(api.resourceMap[path], r)
		}
	}

	return api
}

type mockClient struct {
	// disc is kubernetes discovery client
	disc discovery.DiscoveryInterface
	// dyn is kubernetes dynamic client
	dyn dynamic.Interface
}

func NewMockClient() (api.Client, error) {
	client := fakeclientset.NewSimpleClientset()
	disc, ok := client.Discovery().(*testdiscclient.FakeDiscovery)
	if !ok {
		return nil, fmt.Errorf("couldn't convert Discovery() to *FakeDiscovery")
	}

	dyn := testdynclient.NewSimpleDynamicClient(runtime.NewScheme())

	return &mockClient{
		disc: disc,
		dyn:  dyn,
	}, nil
}

func (m *mockClient) Discover() (api.API, error) {
	return nil, nil
}

func (m *mockClient) Map(api.API) (api.Top, error) {
	return nil, nil
}
