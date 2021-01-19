package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
)

const (
	node1DOTID = "node1ID"
	node2DOTID = "node2ID"
	edgeUID    = "testID"
	weight     = graph.DefaultWeight
)

func TestEdge(t *testing.T) {
	res := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false, api.Options{})

	obj1 := newTestObject(nodeID, nodeName, nodeNs, res, api.Options{})

	n1, err := NewNodeWithDOTID(1, obj1, node1DOTID)
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	obj2 := newTestObject(nodeID, nodeName, nodeNs, res, api.Options{})

	n2, err := NewNodeWithDOTID(2, obj2, node2DOTID)
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	attrs := attrs.New()
	e, err := NewEdgeWithDOTID(edgeUID, n1, n2, weight, entity.Attrs(attrs))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	if uid := e.FromNode().UID(); uid != n1.UID() {
		t.Errorf("expected ID: %s, got: %s", n1.UID(), uid)
	}

	if uid := e.ToNode().UID(); uid != n2.UID() {
		t.Errorf("expected ID: %s, got: %s", n2.UID(), uid)
	}

	if uid := e.DOTID(); uid != edgeUID {
		t.Errorf("expected DOTID: %s, got: %s", edgeUID, uid)
	}

	if w := e.Weight(); big.NewFloat(w).Cmp(big.NewFloat(weight)) != 0 {
		t.Errorf("expected weight %f, got: %f", weight, w)
	}

	re := e.ReversedEdge()

	if re.From().ID() != e.To().ID() {
		t.Errorf("expected from ID: %d, got: %d", e.To().ID(), re.From().ID())
	}

	if re.To().ID() != e.From().ID() {
		t.Errorf("expected to UID: %d, got: %d", e.From().ID(), re.To().ID())
	}

	newDOTID := "DOTID"
	e.SetDOTID(newDOTID)

	if dotID := e.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	if dotAttrs := e.Attributes(); len(dotAttrs) != len(attrs.Attributes()) {
		t.Errorf("expected attributes: %d, got: %d", len(attrs.Attributes()), len(dotAttrs))
	}
}
