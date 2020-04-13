package memory

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
)

func generateAPIObjects() map[string]api.Object {
	a := mock.NewAPI()

	objects := make(map[string]api.Object)

	for _, r := range a.Resources() {
		gv := strings.Join([]string{r.Group(), r.Version()}, "/")

		name := r.Name()
		kind := r.Kind()

		ns := api.NsNan
		if r.Namespaced() {
			ns = mock.Resources[name]["ns"]
		}

		if gvObject, ok := mock.ObjectData[gv]; ok {

			nsKind := strings.Join([]string{ns, kind}, "/")

			if names, ok := gvObject[nsKind]; ok {
				for _, name := range names {
					uid := strings.Join([]string{ns, kind, name}, "/")
					links := make(map[string]api.Relation)
					if rels, ok := mock.ObjectLinks[uid]; ok {
						for obj, rel := range rels {
							links[obj] = mock.NewRelation(rel)
						}
					}
					object := mock.NewObject(name, kind, ns, uid, links)
					objects[uid] = object
				}
			}
		}
	}

	return objects
}

func newTestMemory() (store.Store, error) {
	m := NewStore("testID")

	objects := generateAPIObjects()

	// Store the objects in the memory store
	for _, object := range objects {
		node, err := m.Add(object)
		if err != nil {
			return nil, err
		}

		for _, link := range object.Links() {
			object2 := objects[link.To().String()]

			node2, err := m.Add(object2)
			if err != nil {
				return nil, err
			}

			attrs := store.NewAttributes()
			attrs.Set("relation", link.Relation().String())

			_, err = m.Link(node, node2, store.EntAttrs(attrs))
			if err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func TestNewMemory(t *testing.T) {
	m := NewStore("testID")

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	// NOTE: this test is not needed, but I figured it would be nice
	// to test type-switch into concrete implementation type
	memStore := m.(*Memory)
	expCount := 0
	if nodeCount := memStore.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}
}

func TestAddLinkDelete(t *testing.T) {
	m := NewStore("testID")

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	obj1 := mock.NewObject("foo", "bar", "fobar", "randomid", nil)
	node1, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	node1Obj := node1.Metadata().Get("object")
	node1ApiObj := node1Obj.(api.Object)

	if !reflect.DeepEqual(node1ApiObj, obj1) {
		t.Errorf("expected object: %s, got: %s", obj1, node1ApiObj)
	}

	memStore := m.(*Memory)
	expCount := 1
	if nodeCount := memStore.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	obj2 := mock.NewObject("foo2", "bar2", "fobar", "randomid2", nil)
	node2, err := m.Add(obj2)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	node2Obj := node2.Metadata().Get("object")
	node2ApiObj := node2Obj.(api.Object)

	if !reflect.DeepEqual(node2ApiObj, obj2) {
		t.Errorf("expected object: %s, got: %s", obj2, node2ApiObj)
	}

	expCount = 2
	if nodeCount := memStore.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	edge, err := m.Link(node1, node2)
	if err != nil {
		t.Errorf("failed to link %d to %d: %v", node1.ID(), node2.ID(), err)
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(store.DefaultEdgeWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	if err := m.Delete(edge); err != nil {
		t.Errorf("failed to delete edge: %v", err)
	}

	if edge := m.Edge(node1.ID(), node2.ID()); edge != nil {
		t.Errorf("expected to remove edge between %d-%d, got: %#v", node1.ID(), node2.ID(), edge)
	}

	if err := m.Delete(node1); err != nil {
		t.Errorf("failed to delete node: %v", err)
	}

	if node := m.Node(node1.ID()); node != nil {
		t.Errorf("expected to remove node: %d, got: %#v", node1.ID(), node)
	}

	ent := entity.New()
	if err := m.Delete(ent); err != errors.ErrUnknownEntity {
		t.Errorf("expected: %v, got: %v", errors.ErrUnknownEntity, err)
	}
}

func TestQueryUnknownEntity(t *testing.T) {
	m := NewStore("testID")

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	if _, err := m.Query(); err != errors.ErrUnknownEntity {
		t.Errorf("expected: %v, got: %v", errors.ErrUnknownEntity, err)
	}
}

func TestQueryAllNodes(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	nodes, err := m.Query(query.Entity("node"))
	if err != nil {
		t.Errorf("failed to query all nodes: %v", err)
	}

	if len(nodes) != m.Nodes().Len() {
		t.Errorf("expected node count: %d, got: %d", m.Nodes().Len(), len(nodes))
	}

	for _, nsKinds := range mock.ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]
			for _, name := range names {
				uid := strings.Join([]string{ns, kind, name}, "/")
				nodes, err := m.Query(query.Entity("node"), query.UID(uid))
				if err != nil {
					t.Errorf("error getting node: %s: %v", uid, err)
					continue
				}

				if len(nodes) != 1 {
					t.Errorf("expected single node, got: %d", len(nodes))
					continue
				}

				node := nodes[0]
				object := node.Metadata().Get("object").(api.Object)

				if object.UID().String() != uid {
					t.Errorf("expected node %s, got: %s", uid, object.UID())
				}
			}
		}
	}
}

func TestQueryNodes(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	for _, nsKinds := range mock.ObjectData {
		for nsKind, names := range nsKinds {
			nsplit := strings.Split(nsKind, "/")
			ns, kind := nsplit[0], nsplit[1]

			nodes, err := m.Query(query.Entity("node"), query.Namespace(ns), query.Kind(kind))
			if err != nil {
				t.Errorf("error getting node: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, node := range nodes {
				object := node.Metadata().Get("object").(api.Object)
				if object.Namespace() != ns || object.Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, object.Namespace(), object.Kind())
				}
			}

			for _, name := range names {
				nodes, err := m.Query(query.Entity("node"), query.Namespace(ns), query.Kind(kind), query.Name(name))
				if err != nil {
					t.Errorf("error getting node: %s/%s/%s: %v", ns, kind, name, err)
					continue
				}

				for _, node := range nodes {
					object := node.Metadata().Get("object").(api.Object)
					if object.Namespace() != ns || object.Kind() != kind {
						t.Errorf("expected: %s/%s/%s, got: %s/%s/%s", ns, kind, name,
							object.Namespace(), object.Kind(), object.Name())
					}
				}
			}
		}
	}
}

func TestQueryAllEdges(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	edges, err := m.Query(query.Entity("edge"))
	if err != nil {
		t.Errorf("failed to query edges: %v", err)
	}

	expEdges := 0

	for _, rels := range mock.ObjectLinks {
		expEdges += len(rels)
	}

	if len(edges) != expEdges {
		t.Errorf("expected edge count: %d, got: %d", expEdges, len(edges))
	}
}

func TestQueryAttrEdges(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	attrs := make(map[string]string)

	for _, links := range mock.ObjectLinks {
		for _, relation := range links {
			attrs["relation"] = relation
			edges, err := m.Query(query.Entity("edge"), query.Attrs(attrs))
			if err != nil {
				t.Errorf("failed to query edges with attributes %v: %v", attrs, err)
			}

			for _, edge := range edges {
				for k, v := range attrs {
					if val := edge.Attrs().Get(k); val != v {
						t.Errorf("expected attributes: %v:%v, got: %v:%v", k, v, k, val)
					}
				}
			}
		}
	}
}

func TestSubgraph(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	// NOTE: we are hardcoding this value here
	// as we know that this node UID has 2 neighbouring nodes
	uid := "fooNs/fooKind/foo1"

	nodes, err := m.Query(query.Entity("node"), query.UID(uid))
	if err != nil {
		t.Errorf("failed to find node %s: %v", uid, err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected single node, got: %d", len(nodes))
	}

	node := nodes[0].(store.Node)

	//NOTE: we know the number of expected nodesfrom the moc.ObjectLinks
	testCases := []struct {
		depth int
		exp   int
	}{
		{0, 1},
		{1, 5},
		{100, 6},
	}

	for _, tc := range testCases {
		g, err := m.SubGraph(node, tc.depth)
		if err != nil {
			t.Errorf("failed to query subgraph: %v", err)
		}

		if g.Nodes().Len() != tc.exp {
			t.Errorf("expected subgraph nodes: %d, got: %d", tc.exp, g.Nodes().Len())
		}
	}
}

func TestDOT(t *testing.T) {
	id := "testID"
	m := NewStore(id)

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	dotGraph := m.(store.DOTGraph)
	if dotID := dotGraph.DOTID(); dotID != id {
		t.Errorf("expected DOTID: %s, got: %s", id, dotID)
	}

	graphAttrs, nodeAttrs, edgeAttrs := dotGraph.DOTAttributers()

	memStore := m.(*Memory)

	if !reflect.DeepEqual(graphAttrs, memStore.GraphAttrs) {
		t.Errorf("expected graphtAttrs: %#v, got: %#v", memStore.GraphAttrs, graphAttrs)
	}

	if !reflect.DeepEqual(nodeAttrs, memStore.NodeAttrs) {
		t.Errorf("expected nodeAttrs: %#v, got: %#v", memStore.NodeAttrs, nodeAttrs)
	}

	if !reflect.DeepEqual(edgeAttrs, memStore.EdgeAttrs) {
		t.Errorf("expected edgeAttrs: %#v, got: %#v", memStore.EdgeAttrs, edgeAttrs)
	}

	dot, err := dotGraph.DOT()
	if err != nil {
		t.Errorf("failed to get DOT graph: %v", err)
	}

	if len(dot) == 0 {
		t.Errorf("expected non-empty DOT graph string")
	}
}
