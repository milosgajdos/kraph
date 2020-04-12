package memory

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
)

func TestNewMemory(t *testing.T) {
	m := New("testID")

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
	m := New("testID")

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

func TestDOT(t *testing.T) {
	id := "testID"
	m := New(id)

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	dotStore := m.(store.DOTStore)
	if dotID := dotStore.DOTID(); dotID != id {
		t.Errorf("expected DOTID: %s, got: %s", id, dotID)
	}

	graphAttrs, nodeAttrs, edgeAttrs := dotStore.DOTAttributers()

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

	dot, err := dotStore.DOT()
	if err != nil {
		t.Errorf("failed to get DOT graph: %v", err)
	}

	if len(dot) == 0 {
		t.Errorf("expected non-empty DOT graph string")
	}
}
