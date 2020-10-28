package memory

import (
	"testing"

	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/store/entity"
)

const (
	nodeID   = "testID"
	nodeName = "testName"
)

func TestNode(t *testing.T) {
	attrs := attrs.New()
	attrs.Set("name", nodeName)

	n := entity.NewNode(nodeID, entity.Attrs(attrs))

	node := &Node{
		Node:  n,
		id:    123,
		dotid: nodeName,
	}

	if dotID := node.DOTID(); dotID != nodeName {
		t.Errorf("expected DOTID: %s, got: %s", nodeName, dotID)
	}

	newDOTID := "DOTID"
	node.SetDOTID(newDOTID)

	if dotID := node.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}
}
