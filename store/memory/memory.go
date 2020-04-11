package memory

import (
	"fmt"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
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
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
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
// It never returns error but it might in the future.
func (m *Memory) Add(obj api.Object, opts ...store.Option) (store.Node, error) {
	if id, ok := m.nodes[obj.UID().String()]; ok {
		node := m.WeightedUndirectedGraph.Node(id)
		return node.(store.Node), nil
	}

	id := m.WeightedUndirectedGraph.NewNode().ID()
	name := obj.Kind() + "-" + obj.Name()

	nodeOpts := store.NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	n := entity.NewNode(id, name, store.Meta(nodeOpts.Metadata), store.Attrs(nodeOpts.Attributes))

	n.Metadata().Set("object", obj)

	m.AddNode(n)

	m.nodes[obj.UID().String()] = n.ID()

	return n, nil
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It never returns error but it might in the future.
func (m *Memory) Link(from store.Node, to store.Node, opts ...store.Option) (store.Edge, error) {
	if e := m.Edge(from.ID(), to.ID()); e != nil {
		return e.(store.Edge), nil
	}

	e := entity.NewEdge(from, to, opts...)

	m.SetWeightedEdge(e)

	return e, nil
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
