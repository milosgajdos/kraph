package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
)

var (
	MockAPIResCount               = 9
	odd, even                     = "odd", "even"
	MockAPIOddRes, MockAPIEvenRes = "oddRes", "evenRes"
	MockAPIGroups                 = []string{odd, even}
	MockOddKind, MockEvenKind     = "oddkind", "evenkind"
	MockOddNs, MockEvenNs         = "odd", "even"
	MockAPIMap                    = map[string]map[string]string{
		even: {
			"name":  MockAPIEvenRes,
			"short": "er",
		},
		odd: {
			"name":  MockAPIOddRes,
			"short": "or",
		},
	}
)

func MockAPI() api.API {
	api := &API{
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
		} else {
			r.gv.Group = odd
			r.ar.Name, r.ar.SingularName = MockAPIMap[odd]["name"], MockAPIMap[odd]["short"]
			r.ar.Namespaced = true
		}

		r.gv.Version = fmt.Sprintf("v%d", i)

		api.resources = append(api.resources, r)
		for _, path := range r.Paths() {
			api.resourceMap[path] = append(api.resourceMap[path], r)
		}
	}

	return api
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

	for _, group := range MockAPIGroups {
		name := MockAPIMap[group]["name"]
		for _, res := range a.Resources(query.Name(name)) {
			// create synthetic API objects for given resource map
			for i := 0; i < objCount; i++ {
				ns := MockOddNs
				kind := MockOddKind
				if i%2 == 0 {
					ns = MockEvenNs
					kind = MockEvenKind
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

	// NOTE: we have a map but there are no links between objects, yet

	return top, nil
}

type mockObject struct {
	name  string
	kind  string
	ns    string
	links map[string]map[string]*ObjRef
}

func NewMockObject(name, kind, ns string) api.Object {
	return &mockObject{
		name:  name,
		kind:  kind,
		ns:    ns,
		links: make(map[string]map[string]*ObjRef),
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
	objRef := &ObjRef{
		name: o.Name(),
		kind: o.Kind(),
	}

	key := objRef.name + "/" + objRef.kind
	if m.links[key][r.String()] == nil {
		m.links[key] = make(map[string]*ObjRef)
	}

	if _, ok := m.links[key][r.String()]; !ok {
		m.links[key][r.String()] = objRef
	}

	return nil
}

func (m *mockObject) Links() []api.Link {
	var links []api.Link

	for _, rels := range m.links {
		for rel, obj := range rels {
			link := &Link{
				objRef: obj,
				relation: &Relation{
					r: rel,
				},
			}
			links = append(links, link)
		}
	}

	return links
}
