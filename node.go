package kraph

// NodeOptions configure node
type NodeOptions struct {
	// Attrs are node attributes
	Attrs Attrs
	// Metadata are node metadata
	Metadata Metadata
}

// NodeOption apply node options to node
type NodeOption func(*NodeOptions)

func newNodeOptions(opts ...NodeOption) NodeOptions {
	nodeOpts := NodeOptions{
		Attrs:    make(Attrs),
		Metadata: make(Metadata),
	}

	for _, apply := range opts {
		apply(&nodeOpts)
	}

	if nodeOpts.Attrs == nil {
		nodeOpts.Attrs = make(Attrs)
	}

	if nodeOpts.Metadata == nil {
		nodeOpts.Metadata = make(Metadata)
	}

	return nodeOpts
}

// NodeAttrs sets node attributes
func NodeAttrs(a Attrs) NodeOption {
	return func(o *NodeOptions) {
		o.Attrs = a
	}
}

// NodeMetadata sets node Metadata options
func NodeMetadata(m Metadata) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = m
	}
}

// Node is graph node
type Node struct {
	Attrs
	id       int64
	name     string
	metadata Metadata
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

// Metadata returns node metadata
func (n *Node) Metadata() Metadata {
	return n.metadata
}
