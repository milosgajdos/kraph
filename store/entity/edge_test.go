package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/metadata"
)

var (
	eid        = "edgeID"
	from       = &Node{id: "fooID"}
	to         = &Node{id: "barID"}
	eKey, eVal = "foo", "bar"
)

func newEdgeMeta() *metadata.Metadata {
	meta := metadata.New()
	meta.Set(eKey, eVal)

	return meta
}

func newEdgeAttrs() *attrs.Attrs {
	attrs := attrs.New()
	attrs.Set(aKey, aVal)

	return attrs
}

func TestEdgeID(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata), Attrs(edgeAttrs))

	if e.ID() != eid {
		t.Errorf("expected edge ID: %s, got: %s", eid, e.ID())
	}

	if node := e.From(); node.ID() != from.id {
		t.Errorf("expected from Node: %s, got: %s", from.id, node.ID())
	}

	if node := e.To(); node.ID() != to.id {
		t.Errorf("expected to Node: %s, got: %s", to.id, node.ID())
	}
}

func TestEdgeAttributes(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata), Attrs(edgeAttrs))

	exp := 1

	if count := len(e.Attrs().Keys()); count != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, count)
	}
}

func TestEdgedgeMetadata(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata), Attrs(edgeAttrs))

	if val := e.Metadata().Get(eKey); val != eVal {
		t.Errorf("expected metadata value: %s, got: %s", eVal, val)
	}
}
