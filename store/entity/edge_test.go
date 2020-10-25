package entity

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/attrs"
	"github.com/milosgajdos/kraph/metadata"
)

var (
	eid        = "edgeUID"
	from       = &Node{uid: "fooUID"}
	to         = &Node{uid: "barUID"}
	weight     = 3.0
	eKey, eVal = "foo", "bar"
)

func newEdgeMeta() metadata.Metadata {
	meta := metadata.New()
	meta.Set(eKey, eVal)

	return meta
}

func newEdgeAttrs() attrs.Attrs {
	attrs := attrs.New()
	attrs.Set(aKey, aVal)

	return attrs
}

func TestEdgeUID(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata), Attrs(edgeAttrs))

	if e.UID() != eid {
		t.Errorf("expected edge UID: %s, got: %s", eid, e.UID())
	}

	if node := e.From(); node.UID() != from.uid {
		t.Errorf("expected from Node: %s, got: %s", from.uid, node.UID())
	}

	if node := e.To(); node.UID() != to.uid {
		t.Errorf("expected to Node: %s, got: %s", to.uid, node.UID())
	}
}

func TestEdgeAttributes(t *testing.T) {
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Attrs(edgeAttrs))

	exp := 1

	if count := len(e.Attrs().Keys()); count != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, count)
	}
}

func TestEdgedgeMetadata(t *testing.T) {
	edgeMetadata := newEdgeMeta()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata))

	if val := e.Metadata().Get(eKey); val != eVal {
		t.Errorf("expected metadata value: %s, got: %s", eVal, val)
	}
}

func TestEdgeFromTo(t *testing.T) {
	edgeMetadata := newEdgeMeta()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata))

	fromUid := e.From().UID()
	toUid := e.To().UID()

	if fromUid != from.UID() {
		t.Errorf("expected from UID: %s, got: %s", from.uid, fromUid)
	}

	if toUid != to.UID() {
		t.Errorf("expected to UID: %s, got: %s", to.uid, toUid)
	}
}

func TestEdgeOptions(t *testing.T) {
	edgeMetadata := newEdgeMeta()
	edgeAttrs := newEdgeAttrs()

	e := NewEdge(eid, from, to, Metadata(edgeMetadata), Attrs(edgeAttrs), Weight(weight))

	opts := e.Options()

	if !reflect.DeepEqual(edgeMetadata, opts.Metadata) {
		t.Errorf("expected metadata options: %v, got: %v", edgeMetadata, opts.Metadata)
	}

	if !reflect.DeepEqual(edgeAttrs, opts.Attrs) {
		t.Errorf("expected attributes options: %v, got: %v", edgeAttrs, opts.Attrs)
	}

	if !reflect.DeepEqual(weight, opts.Weight) {
		t.Errorf("expected weight options: %f, got: %f", weight, opts.Weight)
	}
}
