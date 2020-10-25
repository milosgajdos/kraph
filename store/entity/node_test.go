package entity

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/attrs"
	"github.com/milosgajdos/kraph/metadata"
)

var (
	id         = "nodeUID"
	aKey, aVal = "name", "foo"
	mKey, mVal = "foo", "bar"
)

func newNodeMeta() metadata.Metadata {
	meta := metadata.New()
	meta.Set(mKey, mVal)

	return meta
}

func newNodeAttrs() attrs.Attrs {
	attrs := attrs.New()
	attrs.Set(aKey, aVal)

	return attrs
}

func TestNodeUID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Metadata(nodeMetadata), Attrs(nodeAttrs))

	if node.UID() != id {
		t.Errorf("expected node UID: %s, got: %s", id, node.UID())
	}
}

func TestNodeAttributes(t *testing.T) {
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Attrs(nodeAttrs))

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

	node := NewNode(id, Metadata(nodeMetadata))

	if val := node.Metadata().Get(mKey); val != mVal {
		t.Errorf("expected metadata value: %s, got: %s", mVal, val)
	}
}

func TestNodeOptions(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, Attrs(nodeAttrs), Metadata(nodeMetadata))

	opts := node.Options()

	if !reflect.DeepEqual(nodeMetadata, opts.Metadata) {
		t.Errorf("expected metadata options: %v, got: %v", nodeMetadata, opts.Metadata)
	}

	if !reflect.DeepEqual(nodeAttrs, opts.Attrs) {
		t.Errorf("expected attributes options: %v, got: %v", nodeAttrs, opts.Attrs)
	}
}
