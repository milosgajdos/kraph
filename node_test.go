package kraph

import (
	"testing"
)

func TestNode(t *testing.T) {
	var id int64 = 100
	name := "foo"

	node := &Node{
		id:   id,
		name: name,
	}

	if node.ID() != id {
		t.Errorf("invalid id, expected: %d, got: %d", id, node.ID())
	}

	if dotID := node.DOTID(); dotID != name {
		t.Errorf("invalid DOTID, expected: %s, got: %s", name, dotID)
	}

	newID := "bar"
	node.SetDOTID(newID)

	if dotID := node.DOTID(); dotID != newID {
		t.Errorf("invalid DOTID, expected: %s, got: %s", newID, dotID)
	}

	attrsLen := 0
	if nodeattrsLen := len(node.Attributes()); nodeattrsLen != attrsLen {
		t.Errorf("invalid attributes, expected: %d, got: %d", attrsLen, nodeattrsLen)
	}
}
