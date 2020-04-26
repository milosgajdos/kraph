package entity

import "github.com/milosgajdos/kraph/store"

// Node is graph node
type Node struct {
	store.Entity
	id string
}

// NewNode creates a new node and returns it
func NewNode(id string, opts ...store.Option) store.Node {
	return &Node{
		Entity: New(opts...),
		id:     id,
	}
}

// ID returns node ID
func (n *Node) ID() string {
	return n.id
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.Attrs().Get("name")
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Attrs().Set("name", id)
}
