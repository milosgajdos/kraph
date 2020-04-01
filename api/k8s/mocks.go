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
	MockAPIResCount                         = 9
	odd, even                               = "odd", "even"
	MockAPIOddRes, MockAPIEvenRes           = "oddRes", "evenRes"
	MockAPIOddResShort, MockAPIEvenResShort = "or", "er"
	MockAPIGroups                           = []string{odd, even}
	MockOddKind, MockEvenKind               = "oddkind", "evenkind"
	MockOddNs, MockEvenNs                   = "odd", "even"
	// NOTE: the objects are stored under ns/kind/name keys
	MockLinks = map[string]map[string]string{
		"nan/evenkind/evenRes-2": {
			"odd/oddkind/oddRes-1":   "evenodd",
			"nan/evenkind/evenRes-0": "eveneven",
		},
		"odd/oddkind/oddRes-3": {
			"odd/oddkind/oddRes-1":   "oddodd",
			"nan/evenkind/evenRes-0": "oddeven",
		},
	}
	// NOTE: the APIs are stored as api/[name,short]/name
	MockAPIMap = map[string]map[string]string{
		even: {
			"name":  MockAPIEvenRes,
			"short": MockAPIEvenResShort,
		},
		odd: {
			"name":  MockAPIOddRes,
			"short": MockAPIOddResShort,
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

				objName := fmt.Sprintf("%s-%d", res.Name(), i)
				objKind := kind

				var objNs string
				if res.Namespaced() {
					objNs = ns
				}

				if len(objNs) == 0 {
					objNs = NamespaceNan
				}

				objUID := strings.Join([]string{objNs, objKind, objName}, "/")

				if top.index[objNs] == nil {
					top.index[objNs] = make(map[string]map[string]api.Object)
				}

				if top.index[objNs][objKind] == nil {
					top.index[objNs][objKind] = make(map[string]api.Object)
				}

				//fmt.Printf("creating object: %s/%s/%s\n", objNs, objKind, objName)

				obj := NewMockObject(objName, objKind, objNs, objUID)
				top.objects[objUID] = obj
				top.index[objNs][objKind][objName] = obj
			}
		}
	}

	for uid, links := range MockLinks {
		if obj, ok := top.objects[uid]; ok {
			for linkUid, rel := range links {
				if linkObj, lok := top.objects[linkUid]; lok {
					ref := &ObjRef{
						name: linkObj.Name(),
						kind: linkObj.Kind(),
						uid:  linkObj.UID(),
					}
					//fmt.Println("Linking", obj.Name(), "to", ref.Name(), "relation", rel)
					if err := obj.Link(ref, &Relation{r: rel}); err != nil {
						return nil, err
					}
				}
			}
		}
	}

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
