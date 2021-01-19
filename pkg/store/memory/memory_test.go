package memory

import (
	"errors"
	"testing"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
	gm "github.com/milosgajdos/kraph/pkg/graph/memory"
	"github.com/milosgajdos/kraph/pkg/store"
)

func TestNew(t *testing.T) {
	id := "testID"

	m, err := NewStore(id, store.Options{})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	nodes, err := m.Graph().Nodes()
	if err != nil {
		t.Fatalf("failed to get nodes: %v", err)
	}

	expCount := 0
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	if testID := m.ID(); testID != id {
		t.Fatalf("expected id: %s, got: %s", id, testID)
	}

	if g := m.Graph(); m == nil {
		t.Fatalf("expected graph handle, got: %v", g)
	}
}

func TestAddDelete(t *testing.T) {
	m, err := NewStore("testID", store.Options{})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	res := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false, api.Options{})

	node1ID := 1
	node1UID := "foo1UID"
	node1Name := "foo1Name"

	obj := newTestObject(node1UID, node1Name, nodeNs, res, api.Options{})

	n1, err := gm.NewNode(int64(node1ID), obj)
	if err != nil {
		t.Fatalf("failed creating new node: %v", err)
	}

	if err := m.Add(n1, store.NewAddOptions()); err != nil {
		t.Errorf("failed storing node %s: %v", n1.UID(), err)
	}

	node2ID := 2
	node2UID := "foo2UID"
	node2Name := "foo2Name"

	obj2 := newTestObject(node2UID, node2Name, nodeNs, res, api.Options{})

	n2, err := gm.NewNode(int64(node2ID), obj2)
	if err != nil {
		t.Errorf("failed adding node to graph: %v", err)
	}

	if err := m.Add(n2, store.NewAddOptions()); err != nil {
		t.Errorf("failed storing node %s: %v", n2.UID(), err)
	}

	nodes, err := m.Graph().Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount := 2
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	entX, err := entity.NewWithUID("nonExEnt")
	if err != nil {
		t.Fatalf("failed creating entity: %v", err)
	}

	if err := m.Add(entX, store.NewAddOptions()); !errors.Is(err, store.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", store.ErrUnknownEntity, err)
	}

	edge, err := gm.NewEdge(n1, n2, graph.DefaultWeight)
	if err != nil {
		t.Errorf("failed creating edge: %v", err)
	}

	if err := m.Add(edge, store.NewAddOptions()); err != nil {
		t.Errorf("failed storing edge %s: %v", edge.UID(), err)
	}

	edges, err := m.Graph().Edges()
	if err != nil {
		t.Fatalf("failed to get store edges: %v", err)
	}

	expCount = 1
	if edgeCount := len(edges); edgeCount != expCount {
		t.Errorf("expected edges: %d, got: %d", expCount, edgeCount)
	}

	if err := m.Delete(edge, store.DelOptions{}); err != nil {
		t.Errorf("failed deleting edge %s: %v", edge.UID(), err)
	}

	edges, err = m.Graph().Edges()
	if err != nil {
		t.Fatalf("failed to get store edges: %v", err)
	}

	expCount = 0
	if edgeCount := len(edges); edgeCount != expCount {
		t.Errorf("expected edges: %d, got: %d", expCount, edgeCount)
	}

	if err := m.Delete(n2, store.DelOptions{}); err != nil {
		t.Errorf("failed storing node %s: %v", n2.UID(), err)
	}

	nodes, err = m.Graph().Nodes()
	if err != nil {
		t.Fatalf("failed to get store nodes: %v", err)
	}

	expCount = 1
	if nodeCount := len(nodes); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	if err := m.Delete(entX, store.DelOptions{}); !errors.Is(err, store.ErrUnknownEntity) {
		t.Errorf("expected: %v, got: %v", store.ErrUnknownEntity, err)
	}
}
