package memory

import (
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph"
)

type edge struct {
	store.Edge
	from   *node
	to     *node
	weight float64
}

// From returns the from node of the first non-nil edge, or nil.
func (e *edge) From() graph.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns a new Edge with the end point of the edges in the pair swapped
func (e *edge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from

	return e
}

// Weight returns edge weight
func (e *edge) Weight() float64 {
	return e.weight
}
