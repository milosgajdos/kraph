package memory

import (
	"github.com/milosgajdos/kraph/store"
)

// Node is memory store graph Node
type Node struct {
	store.Entity
	id   int64
	name string
}

// ID returns node ID.
// It implements gonum graph.Node interface
func (n *Node) ID() int64 {
	return n.id
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	dotNode, ok := n.Entity.(store.DOTNode)
	if ok {
		return dotNode.DOTID()
	}

	return n.name
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	dotNode, ok := n.Entity.(store.DOTNode)
	if ok {
		dotNode.SetDOTID(id)
	}

	n.name = id
}
