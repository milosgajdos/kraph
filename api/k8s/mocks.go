package k8s

import (
	"fmt"
	"strings"

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
	MockLinks                     = map[string]map[string]string{
		"evenRes-3": {
			"oddRes-1":  "evenodd",
			"evenRes-0": "eveneven",
			"evenRes-1": "eveneven",
		},
		"oddRes-4": {
			"oddRes-1":  "oddodd",
			"oddRes-3":  "oddodd",
			"evenRes-0": "oddeven",
		},
	}
	MockAPIMap = map[string]map[string]string{
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
	top := newTopology()

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

				obj.uid = strings.Join([]string{obj.ns, obj.kind, obj.name}, "-")

				if len(obj.ns) == 0 {
					ns = NamespaceNan
				}

				if top.index[ns] == nil {
					top.index[ns] = make(map[string]map[string]api.Object)
				}

				kind = obj.Kind()
				name = obj.Name()

				if top.index[ns][kind] == nil {
					top.index[ns][kind] = make(map[string]api.Object)
				}

				top.objects[obj.uid] = obj
				top.index[ns][kind][name] = obj
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
	uid   string
	links map[string]map[string]*ObjRef
}

func NewMockObject(name, kind, ns, uid string) api.Object {
	return &mockObject{
		name:  name,
		kind:  kind,
		ns:    ns,
		uid:   uid,
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

func (m *mockObject) UID() string {
	return m.uid
}

func (m *mockObject) Raw() interface{} {
	return m
}

func (m *mockObject) Link(o api.ObjRef, r api.Relation) error {
	objRef := &ObjRef{
		name: o.Name(),
		kind: o.Kind(),
		uid:  o.UID(),
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
