package memory

import "github.com/milosgajdos/kraph/store"

type node struct {
	store.Node
	id   int64
	name string
}

func (n *node) ID() int64 {
	return n.id
}

// DOTID returns the node's DOT ID.
func (n *node) DOTID() string {
	dotNode, ok := n.Node.(store.DOTNode)
	if ok {
		return dotNode.DOTID()
	}

	return n.Node.ID()
}

// SetDOTID sets the node's DOT ID.
func (n *node) SetDOTID(id string) {
	n.name = id
}
