package kraph

import (
	"testing"
)

var (
	id           int64    = 100
	name                  = "foo"
	nKey, nVal            = "foo", "bar"
	nodeMetadata Metadata = map[string]interface{}{
		nKey: nVal,
	}
)

func newNode(name string, id int64, meta Metadata) *Node {
	return &Node{
		Attrs:    make(Attrs),
		id:       id,
		name:     name,
		metadata: meta,
	}
}

func TestNodeID(t *testing.T) {
	node := newNode(name, id, nodeMetadata)

	if node.ID() != id {
		t.Errorf("innValid ID, expected: %d, got: %d", id, node.ID())
	}

}

func TestNodeDOTID(t *testing.T) {
	node := newNode(name, id, nodeMetadata)

	if dotID := node.DOTID(); dotID != name {
		t.Errorf("innValid DOTID, expected: %s, got: %s", name, dotID)
	}

	newID := "bar"
	node.SetDOTID(newID)

	if dotID := node.DOTID(); dotID != newID {
		t.Errorf("innValid DOTID, expected: %s, got: %s", newID, dotID)
	}
}

func TestNodeAttributes(t *testing.T) {
	node := newNode(name, id, nodeMetadata)

	attrsLen := 0
	if nodeattrsLen := len(node.Attributes()); nodeattrsLen != attrsLen {
		t.Errorf("innValid number of attributes, expected: %d, got: %d", attrsLen, nodeattrsLen)
	}
}

func TestNodeMetadata(t *testing.T) {
	node := newNode(name, id, nodeMetadata)

	if meta := node.Metadata(); meta[nKey] != nodeMetadata[nKey] {
		t.Errorf("innValid nodeMetadata, expected: %s, got: %v", nodeMetadata, meta)
	}
}
