package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/entity"
)

const (
	edgeID           = "testID"
	node1ID, node2ID = "node1ID", "node2ID"
)

func TestMemEdge(t *testing.T) {
	n1 := NewNode(1, node1ID, node1ID)
	n2 := NewNode(2, node2ID, node2ID)

	attrs := attrs.New()
	e := NewEdge(edgeID, n1, n2, entity.Attrs(attrs))

	if id := e.From().ID(); id != 1 {
		t.Errorf("expected ID: %d, got: %d", 1, id)
	}

	if id := e.Edge.From().ID(); id != node1ID {
		t.Errorf("expected ID: %s, got: %s", node1ID, id)
	}

	if id := e.To().ID(); id != 2 {
		t.Errorf("expected ID: %d, got: %d", 2, id)
	}

	if id := e.Edge.To().ID(); id != node2ID {
		t.Errorf("expected ID: %s, got: %s", node2ID, id)
	}

	if w := e.Weight(); big.NewFloat(w).Cmp(big.NewFloat(entity.DefaultWeight)) != 0 {
		t.Errorf("expected weight %f, got: %f", entity.DefaultWeight, w)
	}
}
