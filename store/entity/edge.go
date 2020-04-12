package entity

import (
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph"
)

// Edge is graph edge
type Edge struct {
	store.Entity
	from   store.Node
	to     store.Node
	weight float64
}

// NewEdge creates an edge between two nodes and returns it
func NewEdge(from, to store.Node, opts ...store.Option) store.Edge {
	edgeOpts := store.NewOptions()
	for _, apply := range opts {
		apply(&edgeOpts)
	}

	return &Edge{
		Entity: New(opts...),
		from:   from,
		to:     to,
		weight: edgeOpts.Weight,
	}
}

// From returns the from node of the edge
func (e *Edge) From() graph.Node {
	return e.from
}

// To returns the to node of an edge
func (e *Edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns a copy of the edge with reversed nodes
func (e *Edge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from

	return e
}

// Weight returns the edge weight
func (e *Edge) Weight() float64 {
	return e.weight
}
