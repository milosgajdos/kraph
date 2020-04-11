package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/kraph/api/mock"
	"github.com/milosgajdos/kraph/store"
)

func TestMemory(t *testing.T) {
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

func TestAddLink(t *testing.T) {
	m := New("testID")

	if m == nil {
		t.Fatal("failed to create new memory store")
	}

	obj := mock.NewObject("foo", "bar", "fobar", "randomid", nil)
	node1, err := m.Add(obj)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
	}

	memStore := m.(*Memory)
	expCount := 1
	if nodeCount := memStore.Nodes().Len(); nodeCount != expCount {
		t.Errorf("expected nodes: %d, got: %d", expCount, nodeCount)
	}

	obj = mock.NewObject("foo2", "bar2", "fobar", "randomid2", nil)
	node2, err := m.Add(obj)
	if err != nil {
		t.Fatalf("failed adding object to memory store: %v", err)
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
}
