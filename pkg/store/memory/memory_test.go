package memory

import (
	"io/ioutil"
	"math/big"
	"reflect"
	"testing"

	goerr "errors"

	"github.com/ghodss/yaml"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/api/types"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/errors"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/entity"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	objPath = "seeds/objects.yaml"
)

func makeAPIObjects() (map[string]api.Object, error) {
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		return nil, err
	}

	var seedObjects []types.Object
	if err := yaml.Unmarshal(data, &seedObjects); err != nil {
		return nil, err
	}

	objects := make(map[string]api.Object)

	for _, o := range seedObjects {
		res := gen.NewMockResource(
			o.Resource.Name,
			o.Resource.Kind,
			o.Resource.Group,
			o.Resource.Version,
			o.Resource.Namespaced)

		obj := gen.NewMockObject(o.UID, o.Name, o.Namespace, res)

		for _, l := range o.Links {
			obj.Link(uuid.NewFromString(l.To), gen.NewRelation(l.Relation))
		}

		objects[o.UID] = obj
	}

	return objects, nil
}

func newMockObject(uid, name, ns string) api.Object {
	res := gen.NewResource("res", "fooKind", "fooGroup", "v1", true)
	return gen.NewMockObject(uid, name, ns, res)
}

func newTestMemory() (*Memory, error) {
	m, err := NewStore("testID", store.Options{})
	if err != nil {
		return nil, err
	}

	objects, err := makeAPIObjects()
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		node, err := m.Add(object, store.NewAddOptions())
		if err != nil {
			return nil, err
		}

		for _, link := range object.Links() {
			object2 := objects[link.To().String()]

			node2, err := m.Add(object2, store.NewAddOptions())
			if err != nil {
				return nil, err
			}

			attrs := attrs.New()
			attrs.Set("relation", link.Relation().String())

			meta := metadata.New()

			if _, err = m.Link(node, node2, store.LinkOptions{Attrs: attrs, Metadata: meta}); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func TestNewMemory(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get nodes: %v", err)
	}

	expCount := 0
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}
}

func TestAddNode(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	obj := gen.NewMockObject("fooUID", "fooName", "fooNs", nil)

	if _, err := m.Add(obj, store.NewAddOptions()); err != errors.ErrMissingResource {
		t.Errorf("expected error: %v, got: %v", errors.ErrMissingResource, err)
	}

	obj = newMockObject("fooUID", "fooName", "fooNs")

	node1, err := m.Add(obj, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	n, err := m.Node(node1.UID())
	if err != nil {
		t.Fatalf("failed to get node %s: %v", node1.UID(), err)
	}

	if !reflect.DeepEqual(n, node1) {
		t.Errorf("failed getting node %s, got: %v", node1.UID(), n)
	}

	// add the same node again
	nodeX, err := m.Add(obj, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	if !reflect.DeepEqual(node1, nodeX) {
		t.Errorf("expected %s, got %s", node1.UID(), nodeX.UID())
	}
}

func TestGetNode(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	obj := newMockObject("fooUID", "fooName", "fooNs")

	node, err := m.Add(obj, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	nodes, err := m.Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	n, err := m.Node(node.UID())
	if err != nil {
		t.Fatalf("failed to get node %s: %v", node.UID(), err)
	}

	if !reflect.DeepEqual(n, node) {
		t.Errorf("failed getting node %s, got: %v", node.UID(), n)
	}

	if _, err := m.Node(""); err != errors.ErrNodeNotFound {
		t.Errorf("expected %v node, got: %#v", errors.ErrNodeNotFound, err)
	}
}

func TestLink(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	obj1 := newMockObject("fooUID", "fooName", "fooNs")

	node1, err := m.Add(obj1, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	obj2 := newMockObject("foo2UID", "foo2Name", "fooNs")

	node2, err := m.Add(obj2, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	nodeX := entity.NewNode("nonEx")

	if _, err := m.Link(nodeX, node2, store.NewLinkOptions()); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", errors.ErrNodeNotFound, err)
	}

	if _, err := m.Link(node1, nodeX, store.NewLinkOptions()); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected error %s, got: %#v", errors.ErrNodeNotFound, err)
	}

	edge, err := m.Link(node1, node2, store.NewLinkOptions())
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.UID(), node2.UID(), err)
	}

	if w := edge.Weight(); big.NewFloat(w).Cmp(big.NewFloat(store.DefaultWeight)) != 0 {
		t.Errorf("expected non-negative weight")
	}

	edges, err := m.Edges(node1.UID(), node2.UID())
	if err != nil {
		t.Errorf("failed getting edges between %s and %s", node1.UID(), node2.UID())
	}

	if len(edges) == 0 {
		t.Errorf("no edges found between %s and %s", node1.UID(), node2.UID())
	}

	// linking already linked nodes when opts.Line == false
	// must returns the same edge/line as returned previously
	exEdge, err := m.Link(node1, node2, store.NewLinkOptions())
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.UID(), node2.UID(), err)
	}

	if !reflect.DeepEqual(exEdge, edge) {
		t.Errorf("expected %#v, got: %#v", exEdge, edge)
	}

	if _, err := m.Edges("", node2.UID()); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected %v edge, got: %#v", errors.ErrNodeNotFound, err)
	}

	if _, err := m.Edges(node1.UID(), ""); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected %v edge, got: %#v", errors.ErrNodeNotFound, err)
	}
}

func TestDelete(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	obj := newMockObject("fooUID", "fooName", "fooNs")

	node1, err := m.Add(obj, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	obj2 := newMockObject("foo2UID", "foo2Name", "fooNs")

	node2, err := m.Add(obj2, store.NewAddOptions())
	if err != nil {
		t.Fatalf("failed adding object to store: %v", err)
	}

	edge, err := m.Link(node1, node2, store.NewLinkOptions())
	if err != nil {
		t.Errorf("failed to link %s to %s: %v", node1.UID(), node2.UID(), err)
	}

	if err := m.Delete(edge, store.NewDelOptions()); err != nil {
		t.Errorf("failed to delete edge: %v", err)
	}

	edges, err := m.Edges(node1.UID(), node2.UID())
	if err != nil {
		t.Errorf("failed getting edges: %v", err)
	}

	if len(edges) != 0 {
		t.Errorf("expected edges: %d, got: %d", 0, len(edges))
	}

	if err := m.Delete(node1, store.NewDelOptions()); err != nil {
		t.Errorf("failed to delete node: %v", err)
	}

	if _, err := m.Node(node1.UID()); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected %v, got: %v", errors.ErrNodeNotFound, err)
	}

	nodeX := entity.NewNode("nonEx")

	if err := m.Delete(nodeX, store.NewDelOptions()); !goerr.Is(err, errors.ErrNodeNotFound) {
		t.Errorf("expected: %v, got: %v", errors.ErrNodeNotFound, err)
	}

	edgeX := entity.NewEdge("foo", nodeX, nodeX)

	if err := m.Delete(edgeX, store.NewDelOptions()); !goerr.Is(err, errors.ErrEdgeNotFound) {
		t.Errorf("expected: %v, got: %v", errors.ErrNodeNotFound, err)
	}
}

func TestQueryUnknownEntity(t *testing.T) {
	m, err := NewStore("testID", store.NewOptions())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	q := query.Build().Entity("garbage")

	if _, err := m.Query(q); err != errors.ErrInvalidEntity {
		t.Errorf("expected: %v, got: %v", errors.ErrInvalidEntity, err)
	}
}

func TestQueryNodes(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	q := query.Build().MatchAny().Entity(query.Node)

	nodes, err := m.Query(q)
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

	namespaces := make([]string, len(nodes))
	kinds := make([]string, len(nodes))
	names := make([]string, len(nodes))

	for i, n := range nodes {
		o := n.Metadata().Get("object").(api.Object)
		namespaces[i] = o.Namespace()
		kinds[i] = o.Resource().Kind()
		names[i] = o.Name()
	}

	q = query.Build().Entity(query.Node)

	for _, ns := range namespaces {
		q = q.Namespace(ns, query.StringEqFunc(ns))

		nodes, err := m.Query(q)
		if err != nil {
			t.Errorf("error getting namespace %s nodes: %v", ns, err)
			continue
		}

		for _, n := range nodes {
			o := n.Metadata().Get("object").(api.Object)
			if o.Namespace() != ns {
				t.Errorf("expected: namespace %s, got: %s", ns, o.Namespace())
			}
		}

		for _, kind := range kinds {
			q = q.Kind(kind, query.StringEqFunc(kind))

			nodes, err := m.Query(q)
			if err != nil {
				t.Errorf("error getting nodes: %s/%s: %v", ns, kind, err)
				continue
			}

			for _, n := range nodes {
				o := n.Metadata().Get("object").(api.Object)
				if o.Namespace() != ns || o.Resource().Kind() != kind {
					t.Errorf("expected: %s/%s, got: %s/%s", ns, kind, o.Namespace(), o.Resource().Kind())
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

	q := query.Build().MatchAny().Entity(query.Edge)

	edges, err := m.Query(q)
	if err != nil {
		t.Errorf("failed to query edges: %v", err)
	}

	if len(edges) == 0 {
		t.Errorf("expected non-zero edge count")
	}
}

func TestQueryAttrEdges(t *testing.T) {
	m, err := newTestMemory()
	if err != nil {
		t.Fatalf("failed to create new memory store: %v", err)
	}

	q := query.Build().Entity(query.Node)

	nodes, err := m.Query(q)
	if err != nil {
		t.Errorf("failed to query nodes: %v", err)
	}

	relations := make(map[string]bool)

	for _, n := range nodes {
		o := n.Metadata().Get("object").(api.Object)
		for _, l := range o.Links() {
			relations[l.Relation().String()] = true
		}
	}

	a := attrs.New()

	for r := range relations {
		a.Set("relation", r)

		q = query.Build().
			Entity(query.Edge).
			Attrs(a, query.HasAttrsFunc(a))

		edges, err := m.Query(q)
		if err != nil {
			t.Errorf("failed querying edges with attributes %v: %v", a, err)
		}

		for _, edge := range edges {
			for _, k := range a.Keys() {
				v := a.Get(k)
				if val := edge.Attrs().Get(k); val != v {
					t.Errorf("expected attributes: %v:%v, got: %v:%v", k, v, k, val)
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
	// as we know that this node has 2 adjacent nodes
	uid := uuid.NewFromString("fooNs/fooKind/foo1")

	q := query.Build().
		Entity(query.Node).
		UID(uid, query.UIDEqFunc(uid))

	nodes, err := m.Query(q)
	if err != nil {
		t.Errorf("failed to find node %s: %v", uid, err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected single node, got: %d", len(nodes))
	}

	fooNode := NewNode(100, "foo", "bar")

	// subgraph of non-existent node should return error
	if _, err := m.SubGraph(fooNode, 10); err != errors.ErrNodeNotFound {
		t.Errorf("expected: %v, got: %v", errors.ErrNodeNotFound, err)
	}

	node := nodes[0].(store.Node)

	//NOTE: we know the number of expected nodes from seed data
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
	m, err := NewStore(id, store.NewOptions())
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
