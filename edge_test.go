package kraph

import (
	"testing"
)

var (
	weight                = 100.0
	from                  = &Node{id: 1, name: "foo"}
	to                    = &Node{id: 2, name: "bar"}
	eKey, eVal            = "foo", "bar"
	edgeMetadata Metadata = map[string]interface{}{
		eKey: eVal,
	}
)

func newEdge(from, to *Node, weight float64, meta Metadata) *Edge {
	return &Edge{
		Attrs:    make(Attrs),
		from:     from,
		to:       to,
		weight:   weight,
		metadata: meta,
	}
}

func TestEdge(t *testing.T) {
	e := newEdge(from, to, weight, edgeMetadata)

	if node := e.From(); node.ID() != from.id {
		t.Errorf("invalid from node, expected: %d, got: %d", from.id, node.ID())
	}

	if node := e.To(); node.ID() != to.id {
		t.Errorf("invalid to node, expected: %d, got: %d", to.id, node.ID())
	}
}

func TestReversedEdge(t *testing.T) {
	e := newEdge(from, to, weight, edgeMetadata)

	if re := e.ReversedEdge(); re.From().ID() != e.from.ID() || re.To().ID() != e.to.ID() {
		t.Errorf("invalid reversed edge")
	}
}

func TestWeight(t *testing.T) {
	e := newEdge(from, to, weight, edgeMetadata)

	if w := e.Weight(); w != weight {
		t.Errorf("invalid weight, expected: %.2f, got: %.2f", weight, w)
	}
}

func TestEdgeAttributes(t *testing.T) {
	e := newEdge(from, to, weight, edgeMetadata)

	attrsLen := 0
	if edgeAttrsLen := len(e.Attributes()); edgeAttrsLen != attrsLen {
		t.Errorf("invalid attributes, expected: %d, got: %d", attrsLen, edgeAttrsLen)
	}
}

func TestEdgedgeMetadata(t *testing.T) {
	e := newEdge(from, to, weight, edgeMetadata)

	if meta := e.Metadata(); meta[eKey] != edgeMetadata[eKey] {
		t.Errorf("invalid edge Metadata, expected: %s, got: %v", edgeMetadata, meta)
	}
}
