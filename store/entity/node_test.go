package entity

import (
	"testing"

	"github.com/milosgajdos/kraph/store"
)

var (
	id         int64 = 100
	name             = "foo"
	nKey, nVal       = "foo", "bar"
)

func newNodeMeta() store.Metadata {
	meta := store.NewMetadata()
	meta.Set(nKey, nVal)

	return meta
}

func TestNodeID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	node := NewNode(id, name, store.Meta(nodeMetadata))

	if node.ID() != id {
		t.Errorf("expected node ID: %d, got: %d", id, node.ID())
	}
}

func TestNodeDOTID(t *testing.T) {
	nodeMetadata := newNodeMeta()
	node := NewNode(id, name, store.Meta(nodeMetadata))

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
	node := NewNode(id, name, store.Meta(nodeMetadata))

	exp := 0
	if attrsLen := len(node.Attrs().Attributes()); attrsLen != exp {
		t.Errorf("expected attribute count: %d, got: %d", exp, attrsLen)
	}
}

func TestNodeMetadata(t *testing.T) {
	nodeMetadata := newNodeMeta()
	node := NewNode(id, name, store.Meta(nodeMetadata))

	if meta := node.Metadata(); meta.Get(nKey) != nVal {
		t.Errorf("expected metadata value: %s, got: %s", nVal, meta.Get(nKey))
	}
}
