package entity

import (
	"github.com/milosgajdos/kraph/store"
)

// Node implements store.Node
type Node struct {
	id   string
	opts Options
}

// NewNode creates a new node and returns it
func NewNode(id string, opts ...Option) *Node {
	nodeOpts := NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	return &Node{
		id:   id,
		opts: nodeOpts,
	}
}

// UID returns node uid
func (n *Node) UID() string {
	return n.id
}

// Attrs returns node attributes
func (n *Node) Attrs() store.Attrs {
	return n.opts.Attrs
}

// Metadata returns node metadata
func (n *Node) Metadata() store.Metadata {
	return n.opts.Metadata
}

// Options returns node options
func (n Node) Options() Options {
	return n.opts
}
