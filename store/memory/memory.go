package memory

import (
	"fmt"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	// DefaultEdgeWeight defines default edge weight
	DefaultEdgeWeight = 0.0
)

// Memory is in-memory graph store
type Memory struct {
	*simple.WeightedUndirectedGraph
	// id is the store id
	id string
	// nodes maps api.Objects into their node.ID
	nodes map[string]int64
	// Global DOT attributes
	GraphAttrs store.Attributes
	NodeAttrs  store.Attributes
	EdgeAttrs  store.Attributes
}

// New creates new in-memory store and returns it
func New(id string, opts ...store.Option) store.Store {
	o := store.Options{}
	for _, apply := range opts {
		apply(&o)
	}

	if o.GraphAttrs == nil {
		attributes := store.NewAttributes()
		o.GraphAttrs = attributes
	}

	if o.NodeAttrs == nil {
		attributes := store.NewAttributes()
		o.NodeAttrs = attributes
	}

	if o.EdgeAttrs == nil {
		attributes := store.NewAttributes()
		o.EdgeAttrs = attributes
	}

	return &Memory{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		id:                      id,
		nodes:                   make(map[string]int64),
		GraphAttrs:              o.GraphAttrs,
		NodeAttrs:               o.NodeAttrs,
		EdgeAttrs:               o.EdgeAttrs,
	}
}

// Add adds an API object to the in-memory graph as a graph node and returns it
func (m *Memory) Add(obj api.Object, opts ...store.Option) (store.Node, error) {
	return nil, errors.ErrNotImplemented
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
func (m *Memory) Link(from store.Node, to store.Node, opts ...store.Option) (store.Edge, error) {
	return nil, errors.ErrNotImplemented
}

// Query queries the in-memory graph and returns the matched results.
func (m *Memory) Query(q ...query.Option) ([]store.Entity, error) {
	return nil, errors.ErrNotImplemented
}

// DOTID returns the store DOT ID.
func (m *Memory) DOTID() string {
	return m.id
}

// DOTAttributers returns the global DOT graph attributers.
func (m *Memory) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return m.GraphAttrs, m.NodeAttrs, m.EdgeAttrs
}

// DOT returns the GrapViz dot representation of kraph.
func (m *Memory) DOT() (string, error) {
	b, err := dot.Marshal(m, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode kraph into DOT graph: %w", err)
	}

	return string(b), nil
}
