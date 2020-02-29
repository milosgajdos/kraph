package kraph

import "gonum.org/v1/gonum/graph"

// Edge is graph edge
type Edge struct {
	Attrs
	from   *Node
	to     *Node
	weight float64
}

// From returns the from node of the edge
func (e *Edge) From() graph.Node {
	return e.from
}

// To returns the to node of an edge
func (e *Edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns a copy of the edge with reversed edges
func (e *Edge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from

	return e
}

// Weight returns the edge weight
func (e *Edge) Weight() float64 {
	return e.weight
}
