package memory

import (
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph"
)

type Edge struct {
	store.Entity
	from   *Node
	to     *Node
	weight float64
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() graph.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns a new Edge with the end point of the edges in the pair swapped
func (e *Edge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from

	return e
}

// Weight returns edge weight
func (e *Edge) Weight() float64 {
	return e.weight
}
