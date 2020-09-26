package memory

import (
	goerror "errors"
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

func newTestMemory() (*Memory, error) {
	m, err := NewStore("testID")
	if err != nil {
		return nil, err
	}

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

			_, err = m.Link(node.Entity, node2.Entity, store.Attributes(attrs))
			if err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func TestNewMemory(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 0
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}
}

func TestAddNode(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	obj1 := mock.NewObject("foo", "bar", "fobar", "randomid", nil)

	node1, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	n, err := m.Node(node1.Entity.ID())
	if err != nil {
		t.Fatalf("failed to get node %s: %v", node1.Entity.ID(), err)
	}

	if !reflect.DeepEqual(n, node1) {
		t.Errorf("failed getting node %s, got: %v", node1.Entity.ID(), n)
	}

	// add the same node again
	nodeX, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	if !reflect.DeepEqual(node1, nodeX) {
		t.Errorf("expected %s, got %s", node1.Entity.ID(), nodeX.Entity.ID())
	}
}

func TestGetNode(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	obj1 := mock.NewObject("foo", "bar", "fobar", "randomid", nil)

	node1, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	n, err := m.Node(node1.Entity.ID())
	if err != nil {
		t.Fatalf("failed to get node %s: %v", node1.Entity.ID(), err)
	}

	if !reflect.DeepEqual(n, node1) {
		t.Errorf("failed getting node %s, got: %v", node1.Entity.ID(), n)
	}

	if _, err := m.Node(""); err != errors.ErrNodeNotFound {
		t.Errorf("expected %v node, got: %#v", errors.ErrNodeNotFound, err)
	}
}

func TestLink(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	obj1 := mock.NewObject("foo", "bar", "fobar", "randomid", nil)

	node1, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	obj2 := mock.NewObject("foo2", "bar2", "fobar", "randomid2", nil)

	node2, err := m.Add(obj2)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	nodeX := entity.New("nonEx", "foo")

	if _, err := m.Link(nodeX, node2.Entity); !goerror.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", errors.ErrNodeNotFound, err)
	}

	if _, err := m.Link(node1.Entity, nodeX); !goerror.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", errors.ErrNodeNotFound, err)
	}

	edge, err := m.Link(node1.Entity, node2.Entity)
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.Entity.ID(), node2.Entity.ID(), err)
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(store.DefaultEdgeWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	if _, err := m.Edge(node1.Entity.ID(), node2.Entity.ID()); err != nil {
		t.Errorf("failed to find edge between %s and %s", node1.Entity.ID(), node2.Entity.ID())
	}

	exEdge, err := m.Link(node1.Entity, node2.Entity)
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.Entity.ID(), node2.Entity.ID(), err)
	}

	if !reflect.DeepEqual(exEdge, edge) {
		t.Errorf("expected %#v, got: %#v", exEdge, edge)
	}

	if _, err := m.Edge("", node2.Entity.ID()); !goerror.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected %v edge, got: %#v", errors.ErrNodeNotFound, err)
	}

	if _, err := m.Edge(node1.Entity.ID(), ""); !goerror.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected %v edge, got: %#v", errors.ErrNodeNotFound, err)
	}
}

func TestDelete(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
	}

	obj1 := mock.NewObject("foo", "bar", "fobar", "randomid", nil)

	node1, err := m.Add(obj1)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	obj2 := mock.NewObject("foo2", "bar2", "fobar", "randomid2", nil)

	node2, err := m.Add(obj2)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	edge, err := m.Link(node1.Entity, node2.Entity)
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.Entity.ID(), node2.Entity.ID(), err)
	}

	if err := m.Delete(edge); err != nil {
		t.Errorf("failed to delete edge: %v", err)
	}

	if _, err := m.Edge(node1.Entity.ID(), node2.Entity.ID()); err != errors.ErrEdgeNotExist {
		t.Errorf("expected %v, got: %v", errors.ErrEdgeNotExist, err)
	}

	if err := m.Delete(node1.Entity); err != nil {
		t.Errorf("failed to delete node: %v", err)
	}

	if _, err := m.Node(node1.Entity.ID()); err != errors.ErrNodeNotFound {
		t.Errorf("expected %v, got: %v", errors.ErrNodeNotFound, err)
	}

	ent := entity.New("foo", "bar")
	if err := m.Delete(ent); err != errors.ErrUnknownEntity {
		t.Errorf("expected: %v, got: %v", errors.ErrUnknownEntity, err)
	}

	nodeX := entity.New("nonEx", "bar")

	if err := m.Delete(nodeX); err != errors.ErrNodeNotFound {
		t.Errorf("expected: %v, got: %v", errors.ErrNodeNotFound, err)
	}
}

func TestQueryUnknownEntity(t *testing.T) {
	m, err := NewStore("testID")
	if err != nil {
		t.Fatalf("failed to create memory store: %v", err)
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

	storeNodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to fetch store nodes: %v", err)
	}

	if len(nodes) != len(storeNodes) {
		t.Errorf("expected node count: %d, got: %d", len(storeNodes), len(nodes))
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

	// subgraph of non-existent node should return error
	if _, err := m.SubGraph("", 10); err != errors.ErrNodeNotFound {
		t.Errorf("expected: %v, got: %v", errors.ErrNodeNotFound, err)
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
		g, err := m.SubGraph(node.ID(), tc.depth)
		if err != nil {
			t.Errorf("failed to query subgraph: %v", err)
			continue
		}

		storeNodes, err := g.Nodes()
		if err != nil {
			t.Errorf("failed to fetch store nodes: %v", err)
			continue
		}

		if len(storeNodes) != tc.exp {
			t.Errorf("expected subgraph nodes: %d, got: %d", tc.exp, len(storeNodes))
		}
	}
}

func TestDOT(t *testing.T) {
	id := "testID"
	m, err := NewStore(id)
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	if dotID := m.DOTID(); dotID != id {
		t.Errorf("expected DOTID: %s, got: %s", id, dotID)
	}

	dot, err := m.DOT()
	if err != nil {
		t.Errorf("failed to get DOT graph: %v", err)
	}

	if len(dot) == 0 {
		t.Errorf("expected non-empty DOT graph string")
	}
}
