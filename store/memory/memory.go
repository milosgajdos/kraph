package memory

import (
	goerr "errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// Memory is in-memory graph store
type Memory struct {
	*simple.WeightedUndirectedGraph
	// id is the store id
	id string
	// nodes maps api.Objects into their graph Nodes
	nodes map[string]*Node
}

// NewStore creates new in-memory store and returns it
func NewStore(id string, opts ...store.Option) (*Memory, error) {
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &Memory{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		id:                      id,
		nodes:                   make(map[string]*Node),
	}, nil
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (m *Memory) Node(id string) (*Node, error) {
	if node, ok := m.nodes[id]; ok {
		return node, nil
	}

	return nil, errors.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (m *Memory) Nodes() ([]*Node, error) {
	graphNodes := graph.NodesOf(m.WeightedUndirectedGraph.Nodes())

	nodes := make([]*Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise
func (m *Memory) Edge(uid, vid string) (*Edge, error) {
	from, ok := m.nodes[uid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", uid, errors.ErrNodeNotFound)
	}

	to, ok := m.nodes[vid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", vid, errors.ErrNodeNotFound)
	}

	if e := m.WeightedEdge(from.ID(), to.ID()); e != nil {
		return e.(*Edge), nil
	}

	return nil, errors.ErrEdgeNotExist
}

// Add adds an API object to the in-memory graph store and returns it
// It never returns error but it might in the future.
func (m *Memory) Add(obj api.Object, opts ...store.Option) (*Node, error) {
	nodeOpts := store.NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	id := obj.UID().String()
	name := obj.Kind() + "-" + obj.Name()
	nodeOpts.Attributes.Set("name", name)
	entOpts := []store.Option{
		store.Meta(nodeOpts.Metadata),
		store.Attributes(nodeOpts.Attributes),
	}

	e := entity.New(id, name, entOpts...)

	if graphNode, ok := m.nodes[e.ID()]; ok {
		gnode := m.WeightedUndirectedGraph.Node(graphNode.id)
		return gnode.(*Node), nil
	}

	graphNode := m.WeightedUndirectedGraph.NewNode()

	e.Metadata().Set("object", obj)

	node := &Node{
		Entity: e,
		id:     graphNode.ID(),
		name:   e.Name(),
	}

	m.AddNode(node)

	m.nodes[e.ID()] = node

	return node, nil
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It returns error if either of the nodes does not exist in the graph.
func (m *Memory) Link(from store.Node, to store.Node, opts ...store.Option) (*Edge, error) {
	edgeOpts := store.NewOptions()
	for _, apply := range opts {
		apply(&edgeOpts)
	}

	e, err := m.Edge(from.ID(), to.ID())
	if err != nil && err != errors.ErrEdgeNotExist {
		return nil, err
	}

	if e != nil {
		return e, nil
	}

	f, ok := m.nodes[from.ID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	t, ok := m.nodes[to.ID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	id := from.ID() + "-" + to.ID()
	ent := entity.New(id, "edge", opts...)

	edge := &Edge{
		Entity: ent,
		from:   f,
		to:     t,
		weight: edgeOpts.Weight,
	}

	m.SetWeightedEdge(edge)

	return edge, nil
}

// Delete deletes an entity from the memory store
func (m *Memory) Delete(e store.Entity, opts ...store.Option) error {
	switch v := e.(type) {
	case store.Node:
		node, ok := m.nodes[v.ID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.ID(), errors.ErrNodeNotFound)
		}
		m.RemoveNode(node.ID())
		delete(m.nodes, v.ID())
	case store.Edge:
		from, ok := m.nodes[v.From().ID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.From().ID(), errors.ErrNodeNotFound)
		}
		to, ok := m.nodes[v.To().ID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.To().ID(), errors.ErrNodeNotFound)
		}
		m.RemoveEdge(from.ID(), to.ID())
	default:
		return errors.ErrUnknownEntity
	}

	return nil
}

// QueryNode returns all the nodes that match given query.
func (m *Memory) QueryNode(opts ...query.Option) ([]*Node, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	if len(query.UID) > 0 {
		if n, ok := m.nodes[query.UID]; ok {
			return []*Node{n}, nil
		}
	}

	var results []*Node

	visit := func(n graph.Node) {
		node := n.(*Node)
		nodeObj := node.Metadata().Get("object").(api.Object)

		if len(query.Namespace) == 0 || query.Namespace == nodeObj.Namespace() {
			if len(query.Kind) == 0 || query.Kind == nodeObj.Kind() {
				if len(query.Name) == 0 || query.Name == nodeObj.Name() {
					if len(query.Attrs) > 0 {
						for k, v := range query.Attrs {
							if node.Attrs().Get(k) != v {
								return
							}
						}
					}

					// create a deep copy of the matched node
					attrs := store.NewAttributes()
					metadata := store.NewMetadata()

					for _, k := range node.Attrs().Keys() {
						attrs.Set(k, node.Attrs().Get(k))
					}

					for _, k := range node.Metadata().Keys() {
						metadata.Set(k, node.Metadata().Get(k))
					}

					name := node.Entity.Name()
					attrs.Set("name", name)

					e := entity.New(node.Entity.ID(), name, store.Attributes(attrs), store.Meta(metadata))
					node := &Node{Entity: e, id: node.ID(), name: name}

					results = append(results, node)
				}
			}
		}
	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(m.WeightedUndirectedGraph, nil, nil, func(graph.Node) {})

	return results, nil
}

// QueryEdge returns all the edges that match given query
func (m *Memory) QueryEdge(opts ...query.Option) ([]*Edge, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var results []*Edge

	traversed := make(map[int64]map[int64]bool)

	trav := func(e graph.Edge) bool {
		edge := e.(*Edge)

		from := edge.From().(*Node)
		to := edge.To().(*Node)

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
					if edge.Attrs().Get(k) != v {
						return false
					}
				}
			}

			// create a deep copy of the matched edge
			attrs := store.NewAttributes()
			metadata := store.NewMetadata()

			for _, k := range edge.Attrs().Keys() {
				attrs.Set(k, edge.Attrs().Get(k))
			}

			for _, k := range edge.Metadata().Keys() {
				metadata.Set(k, edge.Metadata().Get(k))
			}

			opts := []store.Option{
				store.Weight(edge.Weight()),
				store.Attributes(attrs),
				store.Meta(metadata),
			}

			e := entity.New(edge.Entity.ID(), edge.Entity.Name(), opts...)
			edge := &Edge{Entity: e, from: edge.From().(*Node), to: edge.To().(*Node), weight: edge.Weight()}

			results = append(results, edge)
		}

		return true

	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(m.WeightedUndirectedGraph, nil, nil, func(graph.Node) {})

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
			entities = append(entities, node.Entity)
		}
	case "edge":
		edges, err := m.QueryEdge(q...)
		if err != nil {
			return nil, fmt.Errorf("failed querying edges: %w", err)
		}
		for _, edge := range edges {
			entities = append(entities, edge.Entity)
		}
	default:
		return nil, errors.ErrUnknownEntity
	}

	return entities, nil
}

// SubGraph returns the subgraph of the node up to given depth or returns error
func (m *Memory) SubGraph(id string, depth int) (*Memory, error) {
	rootNode, ok := m.nodes[id]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	s := &Memory{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		id:                      "subgraph" + m.id,
		nodes:                   make(map[string]*Node),
	}

	var sgErr error
	// k2g maps kraph node IDs to subgraph g nodes
	k2g := make(map[int64]*Node)

	visit := func(n graph.Node) {
		vnode := n.(*Node)

		// create a deep copy of the Kraph node
		attrs := store.NewAttributes()
		metadata := store.NewMetadata()

		for _, k := range vnode.Attrs().Keys() {
			attrs.Set(k, vnode.Attrs().Get(k))
		}

		for _, k := range vnode.Metadata().Keys() {
			metadata.Set(k, vnode.Metadata().Get(k))
		}

		obj := vnode.Metadata().Get("object").(api.Object)
		storeNode, err := s.Add(obj, store.Attributes(attrs), store.Meta(metadata))
		if err != nil {
			sgErr = err
			return
		}

		k2g[n.ID()] = storeNode

		// NOTE: this is not very efficient
		// the idea here is we go through newly visited node
		// and check if any of its peer nodes have already
		// been visited (k2g map) and if yes, then wire them
		// to this newly created subgraph node if they
		// have not already been wired to this node (edge is nil)
		nodes := m.From(n.ID())
		for nodes.Next() {
			kraphPeer := nodes.Node().(*Node)
			if to, ok := k2g[kraphPeer.ID()]; ok {
				if _, err := s.Edge(storeNode.Entity.ID(), to.Entity.ID()); goerr.Is(err, errors.ErrEdgeNotExist) {
					// get the original edge from the memory store
					medge, err := m.Edge(vnode.Entity.ID(), kraphPeer.Entity.ID())
					if goerr.Is(err, errors.ErrEdgeNotExist) {
						sgErr = err
						return
					}

					attrs := store.NewAttributes()
					metadata := store.NewMetadata()

					for _, k := range medge.Attrs().Keys() {
						attrs.Set(k, medge.Attrs().Get(k))
					}

					for _, k := range medge.Metadata().Keys() {
						metadata.Set(k, medge.Metadata().Get(k))
					}

					opts := []store.Option{
						store.Weight(medge.Weight()),
						store.Attributes(attrs),
						store.Meta(metadata),
					}

					_, err = s.Link(storeNode.Entity, to.Entity, opts...)
					if err != nil {
						sgErr = err
						return
					}
				}
			}
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	// keep traversing until you reach the requested depth
	_ = bfs.Walk(m.WeightedUndirectedGraph, rootNode, func(n graph.Node, d int) bool {
		if d == depth {
			return true
		}
		return false
	})

	if sgErr != nil {
		return nil, sgErr
	}

	return s, nil
}

// DOTID returns the store DOT ID.
func (m *Memory) DOTID() string {
	return m.id
}

//// DOTAttributers returns the global DOT graph attributers.
//func (m *Memory) DOTAttributers() (graph, node, edge encoding.Attributer) {
//	return m.GraphAttrs, m.NodeAttrs, m.EdgeAttrs
//}

// DOT returns the GrapViz dot representation of kraph.
func (m *Memory) DOT() (string, error) {
	b, err := dot.Marshal(m.WeightedUndirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode into DOT: %w", err)
	}

	return string(b), nil
}
