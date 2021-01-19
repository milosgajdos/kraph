package memory

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// Node is a graph node.
type Node struct {
	entity.Entity
	id    int64
	dotid string
	obj   api.Object
}

// NewNodeWithDOTID creates a new Node with the given DOTID and returns it.
func NewNodeWithDOTID(id int64, obj api.Object, dotid string, opts ...entity.Option) (*Node, error) {
	ent, err := entity.NewWithUID(dotid, opts...)
	if err != nil {
		return nil, err
	}

	return &Node{
		Entity: ent,
		id:     id,
		dotid:  dotid,
		obj:    obj,
	}, nil
}

// NewNode creates a new Node and returns it.
func NewNode(id int64, obj api.Object, opts ...entity.Option) (*Node, error) {
	dotid, err := graph.DOTID(obj)
	if err != nil {
		return nil, err
	}

	eopts := entity.NewOptions()
	for _, apply := range opts {
		apply(&eopts)
	}

	attrs := attrs.NewCopyFrom(eopts.Attrs)
	attrs.Set("dotid", dotid)
	attrs.Set("name", dotid)

	// copy string metadata to node attributes
	for _, k := range obj.Metadata().Keys() {
		if v, ok := obj.Metadata().Get(k).(string); ok {
			attrs.Set(k, v)
		}
	}

	return NewNodeWithDOTID(id, obj, dotid, entity.Attrs(attrs))
}

// Object returns the api.Object this node represents.
func (n Node) Object() api.Object {
	return n.obj
}

// ID returns node ID.
func (n Node) ID() int64 {
	return n.id
}

// DOTID returns GraphViz DOT ID.
func (n Node) DOTID() string {
	return n.dotid
}

// SetDOTID sets GraphViz DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Entity.Attrs().Set("dotid", id)
	n.Entity.Attrs().Set("name", id)
	n.dotid = id
}

// Attributes implements attrs.DOT.
func (n Node) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, len(n.Entity.Attrs().Keys()))

	i := 0
	for _, k := range n.Entity.Attrs().Keys() {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: n.Entity.Attrs().Get(k),
		}
		i++
	}

	return attrs
}
