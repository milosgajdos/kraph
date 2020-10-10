package memory

import (
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph"
)

// TODO: implement graph.Line
// Edge implements graph.WeightedEdge
type Edge struct {
	store.Edge
	from *Node
	to   *Node
	opts entity.Options
}

// NewEdge creates new memory store edge and returns it
func NewEdge(id string, from, to *Node, opts ...entity.Option) *Edge {
	edgeOpts := entity.NewOptions()
	for _, apply := range opts {
		apply(&edgeOpts)
	}

	edge := entity.NewEdge(id, from.Node, to.Node, opts...)

	return &Edge{
		Edge: edge,
		from: from,
		to:   to,
		opts: edgeOpts,
	}
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() graph.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns a new Edge with end points of the edges in the pair swapped
func (e *Edge) ReversedEdge() graph.Edge {
	opts := []entity.Option{
		entity.Attrs(e.opts.Attrs),
		entity.Metadata(e.opts.Metadata),
		entity.Weight(e.opts.Weight),
	}

	edge := entity.NewEdge(e.UID(), e.to.Node, e.from.Node, opts...)

	return &Edge{
		Edge: edge,
		from: e.to,
		to:   e.from,
		opts: e.opts,
	}
}

// Weight returns edge weight
func (e Edge) Weight() float64 {
	return e.opts.Weight
}
