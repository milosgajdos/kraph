package memory

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
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

// Delete deletes an entity from the memory store
func (m *Memory) Delete(e store.Entity, opts ...store.Option) error {
	switch v := e.(type) {
	case store.Node:
		m.RemoveNode(v.ID())
	case store.Edge:
		m.RemoveEdge(v.From().ID(), v.To().ID())
	default:
		return errors.ErrUnknownEntity
	}

	return nil
}

// QueryNode returns all the nodes that match given query.
func (m *Memory) QueryNode(opts ...query.Option) ([]store.Node, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	if len(query.UID) > 0 {
		if id, ok := m.nodes[query.UID]; ok {
			node := m.Node(id)

			return []store.Node{node.(store.Node)}, nil
		}
	}

	var results []store.Node

	visit := func(n graph.Node) {
		node := n.(store.Node)
		nodeObj := node.Metadata().Get("object").(api.Object)

		if len(query.Namespace) == 0 || query.Namespace == nodeObj.Namespace() {
			if len(query.Kind) == 0 || query.Kind == nodeObj.Kind() {
				if len(query.Name) == 0 || query.Name == nodeObj.Name() {
					if len(query.Attrs) > 0 {
						for k, v := range query.Attrs {
							if node.Attributes().Get(k) != v {
								return
							}
						}
					}

					// create a deep copy of the matched node
					attrs := store.NewAttributes()
					metadata := store.NewMetadata()

					for _, k := range node.Attributes().Keys() {
						attrs.Set(k, node.Attributes().Get(k))
					}

					for _, k := range node.Metadata().Keys() {
						metadata.Set(k, node.Metadata().Get(k))
					}

					n := entity.NewNode(node.ID(), node.Name(), store.Attrs(attrs), store.Meta(metadata))

					results = append(results, n)
				}
			}
		}
	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(m, nil, nil, func(graph.Node) {})

	return results, nil
}

// QueryEdge returns all the edges that match given query
func (m *Memory) QueryEdge(opts ...query.Option) ([]store.Edge, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var results []store.Edge

	traversed := make(map[int64]map[int64]bool)

	trav := func(e graph.Edge) bool {
		edge := e.(store.Edge)

		from := edge.From().(store.Node)
		to := edge.To().(store.Node)

		if traversed[from.ID()] == nil {
			traversed[from.ID()] = make(map[int64]bool)
		}

		if traversed[to.ID()] == nil {
			traversed[to.ID()] = make(map[int64]bool)
		}

		if traversed[from.ID()][to.ID()] || traversed[to.ID()][from.ID()] {
			return false
		}

		traversed[from.ID()][to.ID()] = true
		traversed[to.ID()][from.ID()] = true

		if big.NewFloat(query.Weight).Cmp(big.NewFloat(edge.Weight())) == 0 {
			if len(query.Attrs) > 0 {
				for k, v := range query.Attrs {
					if edge.Attributes().Get(k) != v {
						return false
					}
				}
			}

			// create a deep copy of the matched edge
			attrs := store.NewAttributes()
			metadata := store.NewMetadata()

			for _, k := range edge.Attributes().Keys() {
				attrs.Set(k, edge.Attributes().Get(k))
			}

			for _, k := range edge.Metadata().Keys() {
				metadata.Set(k, edge.Metadata().Get(k))
			}

			opts := []store.Option{
				store.Weight(edge.Weight()),
				store.Attrs(attrs),
				store.Meta(metadata),
			}

			e := entity.NewEdge(from, to, opts...)

			results = append(results, e)
		}

		return true

	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(m, nil, nil, func(graph.Node) {})

	return results, nil
}

// Query queries the in-memory graph and returns the matched results.
func (m *Memory) Query(q ...query.Option) ([]store.Entity, error) {
	query := query.NewOptions()
	for _, apply := range q {
		apply(&query)
	}

	var entities []store.Entity

	switch strings.ToLower(query.Entity) {
	case "node":
		nodes, err := m.QueryNode(q...)
		if err != nil {
			return nil, fmt.Errorf("failed querying nodes: %w", err)
		}
		for _, node := range nodes {
			entities = append(entities, node)
		}
	case "edge":
		edges, err := m.QueryEdge(q...)
		if err != nil {
			return nil, fmt.Errorf("failed querying edges: %w", err)
		}
		for _, edge := range edges {
			entities = append(entities, edge)
		}
	default:
		return nil, fmt.Errorf("unknown entity")
	}

	return entities, nil
}

// SubGraph returns the subgraph of the node up to given depth or returns error
func (m *Memory) SubGraph(n store.Node, depth int) (graph.Graph, error) {
	g := simple.NewWeightedUndirectedGraph(0.0, 0.0)

	// k2g maps kraph node IDs to subgraph g nodes
	k2g := make(map[int64]graph.Node)

	visit := func(n graph.Node) {
		node := n.(store.Node)

		// create a deep copy of the Kraph node
		attrs := store.NewAttributes()
		metadata := store.NewMetadata()

		for _, k := range node.Attributes().Keys() {
			attrs.Set(k, node.Attributes().Get(k))
		}

		for _, k := range node.Metadata().Keys() {
			metadata.Set(k, node.Metadata().Get(k))
		}

		gNode := entity.NewNode(g.NewNode().ID(), node.Name(), store.Attrs(attrs), store.Meta(metadata))

		g.AddNode(gNode)
		k2g[n.ID()] = gNode

		// NOTE: this is not very efficient
		// the idea here is we go through newly visited node
		// and check if any of its peer nodes from Kraph have
		// been visited (k2g map) and if yes, then wire them
		// to this newly created subgraph node if they
		// have not already been wired to this node
		nodes := m.From(n.ID())
		for nodes.Next() {
			kraphPeer := nodes.Node()
			if to, ok := k2g[kraphPeer.ID()]; ok {
				if e := g.Edge(gNode.ID(), to.ID()); e == nil {
					edge := m.Edge(n.ID(), kraphPeer.ID())
					kEdge := edge.(store.Edge)

					attrs := store.NewAttributes()
					metadata := store.NewMetadata()

					for _, k := range node.Attributes().Keys() {
						attrs.Set(k, node.Attributes().Get(k))
					}

					for _, k := range node.Metadata().Keys() {
						metadata.Set(k, node.Metadata().Get(k))
					}

					opts := []store.Option{
						store.Weight(kEdge.Weight()),
						store.Attrs(attrs),
						store.Meta(metadata),
					}

					e := entity.NewEdge(gNode, to.(store.Node), opts...)

					g.SetWeightedEdge(e)
				}
			}
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	// keep traversing until you reach the requested depth
	_ = bfs.Walk(m, n, func(n graph.Node, d int) bool {
		if d == depth {
			return true
		}
		return false
	})

	return g, nil
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
		return "", fmt.Errorf("failed to encode into DOT: %w", err)
	}

	return string(b), nil
}
