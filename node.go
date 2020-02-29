package kraph

// Node is graph node
type Node struct {
	Attrs
	// id is node id
	id int64
	// Name names the node
	name string
}

// ID returns node ID
func (n *Node) ID() int64 {
	return n.id
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.name
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.name = id
}
