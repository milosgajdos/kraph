package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store"
)

var (
	id         = "fooID"
	name       = "foo"
	nKey, nVal = "foo", "bar"
)

func newNodeMeta() store.Metadata {
	meta := store.NewMetadata()
	meta.Set(nKey, nVal)

	return meta
}

func newNodeAttrs() store.Attrs {
	attrs := store.NewAttributes()
	attrs.Set("name", name)

	return attrs
}

func TestNodeID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, store.Meta(nodeMetadata), store.EntAttrs(nodeAttrs))

	if node.ID() != id {
		t.Errorf("expected node ID: %s, got: %s", id, node.ID())
	}
}

func TestNodeDOTID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, store.Meta(nodeMetadata), store.EntAttrs(nodeAttrs))

	dotNode := node.(store.DOTNode)

	if dotID := dotNode.DOTID(); dotID != name {
		t.Errorf("expected DOTID: %s, got: %s", name, dotID)
	}

	newID := "bar"
	dotNode.SetDOTID(newID)

	if dotID := dotNode.DOTID(); dotID != newID {
		t.Errorf("expected DOTID: %s, got: %s", newID, dotID)
	}
}

func TestNodeAttributes(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, store.Meta(nodeMetadata), store.EntAttrs(nodeAttrs))

	exp := 1
	if attrsLen := len(node.Attrs().Attributes()); attrsLen != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, attrsLen)
	}
}

func TestNodeMetadata(t *testing.T) {
	nodeMetadata := newNodeMeta()
	nodeAttrs := newNodeAttrs()

	node := NewNode(id, store.Meta(nodeMetadata), store.EntAttrs(nodeAttrs))

	if meta := node.Metadata(); meta.Get(nKey) != nVal {
		t.Errorf("expected metadata value: %s, got: %s", nVal, meta.Get(nKey))
	}
}
