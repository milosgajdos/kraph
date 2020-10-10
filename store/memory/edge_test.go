package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/entity"
)

const (
	edgeID             = "testID"
	node1UID, node2UID = "node1ID", "node2ID"
)

func TestEdge(t *testing.T) {
	n1 := NewNode(1, node1UID, node1UID)
	n2 := NewNode(2, node2UID, node2UID)

	attrs := attrs.New()
	e := NewEdge(edgeID, n1, n2, entity.Attrs(attrs))

	if id := e.From().ID(); id != 1 {
		t.Errorf("expected ID: %d, got: %d", 1, id)
	}

	if id := e.Edge.From().UID(); id != node1UID {
		t.Errorf("expected ID: %s, got: %s", node1UID, id)
	}

	if id := e.To().ID(); id != 2 {
		t.Errorf("expected ID: %d, got: %d", 2, id)
	}

	if id := e.Edge.To().UID(); id != node2UID {
		t.Errorf("expected ID: %s, got: %s", node2UID, id)
	}

	if w := e.Weight(); big.NewFloat(w).Cmp(big.NewFloat(entity.DefaultWeight)) != 0 {
		t.Errorf("expected weight %f, got: %f", entity.DefaultWeight, w)
	}
}
