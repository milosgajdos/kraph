package memory

import (
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node is memory store node
type Node struct {
	store.Node
	id    int64
	dotid string
}

// NewNode creates new memory node and returns it
func NewNode(id int64, uid, dotid string, opts ...entity.Option) *Node {
	node := entity.NewNode(uid, opts...)

	return &Node{
		Node:  node,
		id:    id,
		dotid: dotid,
	}
}

// ID is node ID
func (n *Node) ID() int64 {
	return n.id
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.dotid
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Node.Attrs().Set("name", id)
	n.dotid = id
}

// Attributes implements store.DOTAttrs
func (n *Node) Attributes() []encoding.Attribute {
	if a, ok := n.Attrs().(store.DOTAttrs); ok {
		return a.Attributes()
	}

	attrs := make([]encoding.Attribute, len(n.Attrs().Keys()))

	i := 0
	for _, k := range n.Attrs().Keys() {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: n.Attrs().Get(k),
		}
		i++
	}

	return attrs
}
