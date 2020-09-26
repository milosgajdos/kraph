package memory

import (
	"testing"

	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
)

var (
	nodeID   = "testID"
	nodeName = "testName"
)

func TestMemNode(t *testing.T) {
	attrs := store.NewAttributes()
	attrs.Set("name", nodeName)

	e := entity.New(nodeID, "nodeName", store.Attributes(attrs))

	node := &Node{
		Entity: e,
		id:     123,
		name:   nodeName,
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
