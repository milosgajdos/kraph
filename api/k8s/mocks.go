package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
)

var (
	MockAPIResCount = 9
	odd, even       = "odd", "even"
	mockAPIGroups   = []string{odd, even}
	MockAPIMap      = map[string]map[string]string{
		even: {
			"name":  "evenRes",
			"short": "er",
		},
		odd: {
			"name":  "oddRes",
			"short": "or",
		},
	}
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
			r.gv.Group = even
			r.ar.Name, r.ar.SingularName = MockAPIMap[even]["name"], MockAPIMap[even]["short"]
			r.ar.Namespaced = true
		} else {
			r.gv.Group = odd
			r.ar.Name, r.ar.SingularName = MockAPIMap[odd]["name"], MockAPIMap[odd]["short"]
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

type mockClient struct{}

func NewMockClient() (api.Client, error) {
	return &mockClient{}, nil
}

func (m *mockClient) Discover() (api.API, error) {
	return MockAPI(), nil
}

func (m *mockClient) Map(a api.API) (api.Top, error) {
	top := make(Top)

	objCount := 5

	for _, group := range mockAPIGroups {
		name := MockAPIMap[group]["name"]
		for _, res := range a.Resources(query.Name(name)) {
			// create synthetic API objects for given resource map
			for i := 0; i < objCount; i++ {
				ns := "odd"
				kind := "oddkind"
				if i%2 == 0 {
					ns = "even"
					kind = "evenkind"
				}

				obj := &mockObject{
					name: fmt.Sprintf("%s-%d", res.Name(), i),
					kind: kind,
				}

				if res.Namespaced() {
					obj.ns = ns
				}

				if len(obj.ns) == 0 {
					ns = NamespaceNan
				}

				if top[ns] == nil {
					top[ns] = make(map[string]map[string]api.Object)
				}

				kind = obj.Kind()
				name = obj.Name()

				if top[ns][kind] == nil {
					top[ns][kind] = make(map[string]api.Object)
				}

				top[ns][kind][name] = obj
			}
		}
	}

	// TODO: we have a map but we don't have any links between objects, yet

	return top, nil
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

func (m *mockObject) Link(o api.ObjRef, r api.Relation) error {
	objRef := ObjRef{
		name: o.Name(),
		kind: o.Kind(),
	}

	link := &Link{
		objRef: objRef,
		relation: &Relation{
			r: r.String(),
		},
	}

	m.links = append(m.links, link)

	return nil
}

func (m *mockObject) Links() []api.Link {
	lx := make([]api.Link, len(m.links))

	for i, l := range m.links {
		lx[i] = l
	}

	return lx
}
