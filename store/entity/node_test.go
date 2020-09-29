package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/metadata"
)

var (
	id         = "nodeID"
	aKey, aVal = "name", "foo"
	mKey, mVal = "foo", "bar"
)

func newNodeMeta() *metadata.Metadata {
	meta := metadata.New()
	meta.Set(mKey, mVal)

	return meta
}

func newNodeAttrs() *attrs.Attrs {
	attrs := attrs.New()
	attrs.Set(aKey, aVal)

	return attrs
}

func TestNodeID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Metadata(nodeMetadata), Attrs(nodeAttrs))

	if node.ID() != id {
		t.Errorf("expected node ID: %s, got: %s", id, node.ID())
	}
}

func TestNodeAttributes(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Metadata(nodeMetadata), Attrs(nodeAttrs))

	exp := 1
	if attrsLen := len(node.Attrs().Keys()); attrsLen != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, attrsLen)
	}

	if val := node.Attrs().Get(aKey); val != aVal {
		t.Errorf("expected attr val: %s, got: %s", aVal, val)
	}
}

func TestNodeMetadata(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Metadata(nodeMetadata), Attrs(nodeAttrs))

	if val := node.Metadata().Get(mKey); val != mVal {
		t.Errorf("expected metadata value: %s, got: %s", mVal, val)
	}
}
