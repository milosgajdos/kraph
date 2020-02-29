package kraph

import "testing"

func TestEdge(t *testing.T) {
	weight := 100.0

	from := &Node{id: 1, name: "foo"}
	to := &Node{id: 2, name: "bar"}

	e := &Edge{
		from:   from,
		to:     to,
		weight: weight,
	}

	if node := e.From(); node.ID() != from.id {
		t.Errorf("invalid from node, expected: %d, got: %d", from.id, node.ID())
	}

	if node := e.To(); node.ID() != to.id {
		t.Errorf("invalid to node, expected: %d, got: %d", to.id, node.ID())
	}

	if re := e.ReversedEdge(); re.From().ID() != e.from.ID() || re.To().ID() != e.to.ID() {
		t.Errorf("invalid reversed edge")
	}

	if w := e.Weight(); w != weight {
		t.Errorf("invalid weight, expected: %.2f, got: %.2f", weight, w)
	}

	attrsLen := 0
	if edgeAttrsLen := len(e.Attributes()); edgeAttrsLen != attrsLen {
		t.Errorf("invalid attributes, expected: %d, got: %d", attrsLen, edgeAttrsLen)
	}
}
