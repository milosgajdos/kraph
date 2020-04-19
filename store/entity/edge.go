package entity

import (
	"github.com/milosgajdos/kraph/store"
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
func (e *Edge) From() store.Node {
	return e.from
}

// To returns the to node of an edge
func (e *Edge) To() store.Node {
	return e.to
}

// Weight returns the edge weight
func (e *Edge) Weight() float64 {
	return e.weight
}
