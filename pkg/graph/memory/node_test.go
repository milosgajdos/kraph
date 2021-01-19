package memory

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
)

func TestNode(t *testing.T) {
	obj := newTestObject(nodeID, nodeName, nodeNs, nil, api.Options{})

	node, err := NewNode(nodeGID, obj)
	if err == nil {
		t.Fatalf("expected error %v, got: %v", graph.ErrMissingResource, err)
	}

	res := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false, api.Options{})
	obj = newTestObject(nodeID, nodeName, nodeNs, res, api.Options{})

	dotid, err := graph.DOTID(obj)
	if err != nil {
		t.Fatalf("failed to build DOTID: %v", err)
	}

	attrs := attrs.New()
	attrs.Set("nodename", nodeName)

	node, err = NewNode(nodeGID, obj, entity.Attrs(attrs))
	if err != nil {
		t.Fatalf("failed to create new node from API object: %v", err)
	}

	if id := node.ID(); id != nodeGID {
		t.Errorf("expected ID: %d, got: %d", nodeGID, id)
	}

	if nodeObj := node.Object(); !reflect.DeepEqual(nodeObj, obj) {
		t.Errorf("invalid api.Object for node: %s", node.UID())
	}

	if dotID := node.DOTID(); dotID != dotid {
		t.Errorf("expected DOTID: %s, got: %s", dotid, dotID)
	}

	// NOTE: by default we will get the following attributes:
	// * dotid
	// * name
	// We add "nodename" attribute which leaves us with 3 attributes altogether.
	if dotAttrs := node.Attributes(); len(dotAttrs) != 3 {
		t.Errorf("expected %d attributes, got: %d", 3, len(dotAttrs))
	}

	newDOTID := "DOTID"
	node.SetDOTID(newDOTID)

	if dotID := node.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}
}

func TestNodeWithDOTID(t *testing.T) {
	res := newTestResource(nodeResName, nodeResGroup, nodeResVersion, nodeResKind, false, api.Options{})
	obj := newTestObject(nodeID, nodeName, nodeNs, res, api.Options{})

	attrs := attrs.New()
	attrs.Set("name", nodeName)

	node, err := NewNodeWithDOTID(nodeGID, obj, nodeName, entity.Attrs(attrs))
	if err != nil {
		t.Errorf("failed to create new node: %v", err)
		return
	}

	if id := node.ID(); id != nodeGID {
		t.Errorf("expected ID: %d, got: %d", nodeGID, id)
	}

	if dotID := node.DOTID(); dotID != nodeName {
		t.Errorf("expected DOTID: %s, got: %s", nodeName, dotID)
	}

	newDOTID := "DOTID"
	node.SetDOTID(newDOTID)

	if dotID := node.DOTID(); dotID != newDOTID {
		t.Errorf("expected DOTID: %s, got: %s", newDOTID, dotID)
	}

	if dotAttrs := node.Attributes(); len(dotAttrs) != len(attrs.Attributes()) {
		t.Errorf("expected attributes: %d, got: %d", len(attrs.Attributes()), len(dotAttrs))
	}
}
