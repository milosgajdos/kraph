package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store"
)

var (
	weight     = 100.0
	from       = &Node{Entity: New(), id: "fooID"}
	to         = &Node{Entity: New(), id: "barID"}
	eKey, eVal = "foo", "bar"
)

func newEdgeMeta() store.Metadata {
	meta := store.NewMetadata()
	meta.Set(eKey, eVal)

	return meta
}

func TestEdge(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	e := NewEdge(from, to, store.Weight(weight), store.Meta(edgeMetadata))

	if node := e.From(); node.ID() != from.id {
		t.Errorf("expected from Node: %s, got: %s", from.id, node.ID())
	}

	if node := e.To(); node.ID() != to.id {
		t.Errorf("expected to Node: %s, got: %s", to.id, node.ID())
	}
}

func TestWeight(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	e := NewEdge(from, to, store.Weight(weight), store.Meta(edgeMetadata))

	if w := e.Weight(); w != weight {
		t.Errorf("expected weight: %.2f, got: %.2f", weight, w)
	}
}

func TestEdgeAttributes(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	e := NewEdge(from, to, store.Weight(weight), store.Meta(edgeMetadata))

	exp := 0
	if attrsLen := len(e.Attrs().Attributes()); attrsLen != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, attrsLen)
	}
}

func TestEdgedgeMetadata(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	e := NewEdge(from, to, store.Weight(weight), store.Meta(edgeMetadata))

	if meta := e.Metadata(); meta.Get(eKey) != eVal {
		t.Errorf("expected metadata value: %s, got: %s", eVal, meta.Get(eKey))
	}
}
