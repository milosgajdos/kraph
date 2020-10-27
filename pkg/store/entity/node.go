package entity

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node implements store.Node
type Node struct {
	uid  string
	opts Options
}

// NewNode creates a new node and returns it
func NewNode(uid string, opts ...Option) *Node {
	nodeOpts := NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	return &Node{
		uid:  uid,
		opts: nodeOpts,
	}
}

// UID returns node uid
func (n *Node) UID() string {
	return n.uid
}

// Attrs returns node attributes
func (n *Node) Attrs() attrs.Attrs {
	return n.opts.Attrs
}

// Attributes returns attributes as a slice of encoding.Attribute
func (n *Node) Attributes() []encoding.Attribute {
	return attrs.DOTAttrs(n.opts.Attrs)
}

// Metadata returns node metadata
func (n *Node) Metadata() metadata.Metadata {
	return n.opts.Metadata
}

// Options returns node options
func (n Node) Options() Options {
	return n.opts
}
