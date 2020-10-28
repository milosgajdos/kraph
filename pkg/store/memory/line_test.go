package memory

import (
	"math/big"
	"testing"

	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/store/entity"
)

const (
	lineID, lineUID    = 10, "testID"
	node1UID, node2UID = "node1ID", "node2ID"
)

func TestLine(t *testing.T) {
	n1 := NewNode(1, node1UID, node1UID)
	n2 := NewNode(2, node2UID, node2UID)

	attrs := attrs.New()
	l := NewLine(lineID, lineUID, lineUID, n1, n2, entity.Attrs(attrs))

	if id := l.ID(); id != lineID {
		t.Errorf("expected ID: %d, got: %d", lineID, id)
	}

	if id := l.Edge.From().UID(); id != node1UID {
		t.Errorf("expected ID: %s, got: %s", node1UID, id)
	}

	if id := l.To().ID(); id != 2 {
		t.Errorf("expected ID: %d, got: %d", 2, id)
	}

	if id := l.Edge.To().UID(); id != node2UID {
		t.Errorf("expected ID: %s, got: %s", node2UID, id)
	}

	if id := l.DOTID(); id != lineUID {
		t.Errorf("expected DOTID: %s, got: %s", lineUID, id)
	}

	if w := l.Weight(); big.NewFloat(w).Cmp(big.NewFloat(entity.DefaultWeight)) != 0 {
		t.Errorf("expected weight %f, got: %f", entity.DefaultWeight, w)
	}

	rl := l.ReversedLine()

	if rl.From().ID() != l.To().ID() {
		t.Errorf("expected from ID: %d, got: %d", l.To().ID(), rl.From().ID())
	}

	if rl.To().ID() != l.From().ID() {
		t.Errorf("expected to ID: %d, got: %d", l.From().ID(), rl.To().ID())
	}
}
