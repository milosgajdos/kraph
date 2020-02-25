package kraph

import (
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func TestNode(t *testing.T) {
	g := simple.NewWeightedUndirectedGraph(0.0, 0.0)
	name := "foo"
	node := &Node{
		Node: g.NewNode(),
		Name: name,
	}

	if dotID := node.DOTID(); dotID != name {
		t.Errorf("expected: %s, go: %s", name, dotID)
	}

	id := "bar"
	node.SetDOTID(id)

	if dotID := node.DOTID(); dotID != id {
		t.Errorf("expected: %s, go: %s", id, dotID)
	}
}
