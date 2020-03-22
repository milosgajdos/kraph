package kraph

import (
	"gonum.org/v1/gonum/graph"
)

// EdgeOptions are edge options
type EdgeOptions struct {
	// Attrs are edge attributes
	Attrs Attrs
	// Metadata are edge metadata
	Metadata Metadata
	// Weight is edge weight
	Weight float64
}

// EdgeOption applies edge options
type EdgeOption func(*EdgeOptions)

func newEdgeOptions(opts ...EdgeOption) EdgeOptions {
	edgeOpts := EdgeOptions{
		Weight:   DefaultWeight,
		Attrs:    make(Attrs),
		Metadata: make(Metadata),
	}

	for _, apply := range opts {
		apply(&edgeOpts)
	}

	if edgeOpts.Attrs == nil {
		edgeOpts.Attrs = make(Attrs)
	}

	if edgeOpts.Metadata == nil {
		edgeOpts.Metadata = make(Metadata)
	}

	return edgeOpts
}

// EdgeAttrs sets edge attributes
func EdgeAttrs(a Attrs) EdgeOption {
	return func(o *EdgeOptions) {
		o.Attrs = a
	}
}

// EdgeMetadata sets edge metadata
func EdgeMetadata(m Metadata) EdgeOption {
	return func(o *EdgeOptions) {
		o.Metadata = m
	}
}

// Weight set edge weight
func Weight(w float64) EdgeOption {
	return func(o *EdgeOptions) {
		o.Weight = w
	}
}

// Edge is graph edge
type Edge struct {
	Attrs
	from     *Node
	to       *Node
	weight   float64
	metadata Metadata
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

// Metadata returns edge metadata
func (e *Edge) Metadata() Metadata {
	return e.metadata
}
