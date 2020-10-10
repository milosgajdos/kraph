package entity

import (
	"github.com/milosgajdos/kraph/store"
)

// Edge implements store.Edge
type Edge struct {
	id   string
	from store.Node
	to   store.Node
	opts Options
}

// NewEdge creates an edge between two nodes and returns it
func NewEdge(id string, from, to store.Node, opts ...Option) *Edge {
	edgeOpts := NewOptions()
	for _, apply := range opts {
		apply(&edgeOpts)
	}

	return &Edge{
		id:   id,
		from: from,
		to:   to,
		opts: edgeOpts,
	}
}

// UID returns edge uid
func (e *Edge) UID() string {
	return e.id
}

// Attrs returns edge attributes
func (e *Edge) Attrs() store.Attrs {
	return e.opts.Attrs
}

// Metadata reutnrs edge metadata
func (e *Edge) Metadata() store.Metadata {
	return e.opts.Metadata
}

// From returns the from node of the edge
func (e *Edge) From() store.Node {
	return e.from
}

// To returns the to node of an edge
func (e *Edge) To() store.Node {
	return e.to
}

// Weight returns edge weight
func (e *Edge) Weight() float64 {
	return e.opts.Weight
}

// Options return options
func (e Edge) Options() Options {
	return e.opts
}
