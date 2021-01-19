package memory

import (
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
	"github.com/milosgajdos/kraph/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// Edge implements graph.WeightedEdge
type Edge struct {
	entity.Entity
	from   *Node
	to     *Node
	dotid  string
	weight float64
}

// NewEdgeWithDOTID creates a new edge with the given dotid and returns it.
func NewEdgeWithDOTID(dotid string, from, to *Node, weight float64, opts ...entity.Option) (*Edge, error) {
	ent, err := entity.NewWithUID(dotid, opts...)
	if err != nil {
		return nil, err
	}

	return &Edge{
		Entity: ent,
		from:   from,
		to:     to,
		dotid:  dotid,
		weight: weight,
	}, nil
}

// NewEdge creates a new edge and returns it.
// NewEdge sets the edge DOTID to uid.
func NewEdge(from, to *Node, weight float64, opts ...entity.Option) (*Edge, error) {
	uid := uuid.New().String()

	return NewEdgeWithDOTID(uid, from, to, weight, opts...)
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() gngraph.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() gngraph.Node {
	return e.to
}

// FromNode returns the from node of the first non-nil edge, or nil.
func (e *Edge) FromNode() graph.Node {
	return e.from
}

// ToNode returns the to node of the first non-nil edge, or nil.
func (e *Edge) ToNode() graph.Node {
	return e.to
}

// Weight returns edge weight
func (e Edge) Weight() float64 {
	return e.weight
}

// ReversedEdge returns a new line with end points of the pair swapped
func (e *Edge) ReversedEdge() gngraph.Edge {
	return &Edge{
		Entity: e.Entity,
		from:   e.to,
		to:     e.from,
		dotid:  e.dotid,
		weight: e.weight,
	}
}

// DOTID returns the edge's DOT ID.
func (e *Edge) DOTID() string {
	return e.dotid
}

// SetDOTID sets the edge's DOT ID.
func (e *Edge) SetDOTID(id string) {
	e.Attrs().Set("dotid", id)
	e.dotid = id
}

// Attributes implements store.DOTAttrs
func (e *Edge) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, len(e.Attrs().Keys()))

	i := 0
	for _, k := range e.Attrs().Keys() {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: e.Attrs().Get(k),
		}
		i++
	}

	return attrs
}
