package kraph

import "gonum.org/v1/gonum/graph"

// Node is graph node
type Node struct {
	graph.Node
	Attrs
	// Name names the node
	Name string
}

// ID returns node ID
func (n *Node) ID() int64 {
	return n.Node.ID()
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.Name
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Name = id
}
