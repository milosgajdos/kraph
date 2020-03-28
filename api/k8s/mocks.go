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

type mockAPI struct {
	*API
}

func MockAPI() api.API {
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

	return &mockAPI{
		API: &api,
	}
}

type mockTop struct {
	Top
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
	return MockAPI(), nil
}

func (m *mockClient) Map(api.API) (api.Top, error) {
	top := make(Top)

	return &mockTop{
		Top: top,
	}, nil
}

type mockObject struct {
	name  string
	kind  string
	ns    string
	links []*Link
}

func NewMockObject(name, kind, ns string, links ...*Link) api.Object {
	oLinks := make([]*Link, len(links))

	for i, o := range links {
		oLinks[i] = o
	}

	return &mockObject{
		name:  name,
		kind:  kind,
		ns:    ns,
		links: oLinks,
	}
}

func (m *mockObject) Name() string {
	return m.name
}

func (m *mockObject) Kind() string {
	return m.kind
}

func (m *mockObject) Namespace() string {
	return m.ns
}

func (m *mockObject) Raw() interface{} {
	return m
}

func (m *mockObject) Links() []api.Link {
	lx := make([]api.Link, len(m.links))

	for i, l := range m.links {
		lx[i] = l
	}

	return lx
}
