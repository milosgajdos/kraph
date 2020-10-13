package memory

import (
	"fmt"
	"strings"

	goerr "errors"

	"github.com/google/uuid"
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/entity"
	"github.com/milosgajdos/kraph/store/metadata"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// Memory is in-memory graph store
type Memory struct {
	// TODO: make graph configurable
	*simple.WeightedUndirectedGraph
	// id is the store id
	id string
	// nodes maps api.Objects into their graph Nodes
	nodes map[string]*Node
	// edges maps links between api.Object nodes to their graph Edge
	edges map[string]*Edge
	// options are store options
	opts store.Options
}

// NewStore creates new in-memory store and returns it
func NewStore(id string, opts store.Options) (*Memory, error) {
	return &Memory{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(1.0, 1.0),
		id:                      id,
		nodes:                   make(map[string]*Node),
		edges:                   make(map[string]*Edge),
		opts:                    opts,
	}, nil
}

// Add adds an API object to the in-memory graph store and returns its entity.
func (m *Memory) Add(obj api.Object, opts store.AddOptions) (store.Entity, error) {
	uid := obj.UID().String()

	if node, ok := m.nodes[uid]; ok {
		return node, nil
	}

	// make a copy of store attributes and metadata
	attrs := attrs.New()
	metadata := metadata.New()

	if opts.Attrs != nil {
		for _, k := range opts.Attrs.Keys() {
			attrs.Set(k, opts.Attrs.Get(k))
		}
	}

	if opts.Metadata != nil {
		for _, k := range opts.Metadata.Keys() {
			metadata.Set(k, opts.Metadata.Get(k))
		}
	}

	ns := obj.Namespace()
	// TODO: figure out what to do here
	if obj.Namespace() == "" || obj.Namespace() == api.NsNan {
		ns = "global"
	}

	dotid := strings.Join([]string{ns, obj.Kind(), obj.Name()}, "/")
	attrs.Set("name", dotid)

	entOpts := []entity.Option{
		entity.Metadata(metadata),
		entity.Attrs(attrs),
	}

	g := m.WeightedUndirectedGraph.NewNode()

	node := NewNode(g.ID(), uid, dotid, entOpts...)

	node.Metadata().Set("object", obj)

	m.WeightedUndirectedGraph.AddNode(node)

	m.nodes[uid] = node

	return node, nil
}

// Delete deletes an entity from the memory store
func (m *Memory) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case store.Edge:
		from, ok := m.nodes[v.From().UID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.From().UID(), errors.ErrNodeNotFound)
		}
		to, ok := m.nodes[v.To().UID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.To().UID(), errors.ErrNodeNotFound)
		}
		m.RemoveEdge(from.ID(), to.ID())
		delete(m.nodes, v.UID())
	case store.Node:
		node, ok := m.nodes[v.UID()]
		if !ok {
			return fmt.Errorf("%s: %w", v.UID(), errors.ErrNodeNotFound)
		}
		m.WeightedUndirectedGraph.RemoveNode(node.ID())
		delete(m.nodes, v.UID())
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
					attrs := attrs.New()
					metadata := metadata.New()

					for _, k := range node.Attrs().Keys() {
						attrs.Set(k, node.Attrs().Get(k))
					}

					for _, k := range node.Metadata().Keys() {
						metadata.Set(k, node.Metadata().Get(k))
					}

					dotid := strings.Join([]string{nodeObj.Namespace(), nodeObj.Kind(), nodeObj.Name()}, "/")
					attrs.Set("name", dotid)

					entOpts := []entity.Option{
						entity.Metadata(metadata),
						entity.Attrs(attrs),
					}

					n := NewNode(node.ID(), node.UID(), dotid, entOpts...)

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

		if query.Weight != nil {
			if !query.Weight(edge.Weight()) {
				return false
			}
		}

		if len(query.Attrs) > 0 {
			for k, v := range query.Attrs {
				if edge.Attrs().Get(k) != v {
					return false
				}
			}
		}

		// create a deep copy of the matched edge
		attrs := attrs.New()
		metadata := metadata.New()

		for _, k := range edge.Attrs().Keys() {
			attrs.Set(k, edge.Attrs().Get(k))
		}

		for _, k := range edge.Metadata().Keys() {
			metadata.Set(k, edge.Metadata().Get(k))
		}

		opts := []entity.Option{
			entity.Attrs(attrs),
			entity.Metadata(metadata),
			entity.Weight(edge.Weight()),
		}

		ent := NewEdge(edge.UID(), from, to, opts...)

		results = append(results, ent)

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
			entities = append(entities, node.Node)
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
		return nil, errors.ErrUnknownEntity
	}

	return entities, nil
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (m *Memory) Node(id string) (store.Node, error) {
	if node, ok := m.nodes[id]; ok {
		return node, nil
	}

	return nil, errors.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (m *Memory) Nodes() ([]store.Node, error) {
	graphNodes := graph.NodesOf(m.WeightedUndirectedGraph.Nodes())

	nodes := make([]store.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise
func (m *Memory) Edge(uid, vid string) (store.Edge, error) {
	from, ok := m.nodes[uid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", uid, errors.ErrNodeNotFound)
	}

	to, ok := m.nodes[vid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", vid, errors.ErrNodeNotFound)
	}

	if e := m.WeightedEdge(from.ID(), to.ID()); e != nil {
		return e.(*Edge).Edge, nil
	}

	return nil, errors.ErrEdgeNotExist
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It returns error if either of the nodes does not exist in the graph.
func (m *Memory) Link(from store.Node, to store.Node, opts store.LinkOptions) (store.Edge, error) {
	e, err := m.Edge(from.UID(), to.UID())
	if err != nil && err != errors.ErrEdgeNotExist {
		return nil, err
	}

	if e != nil {
		return e, nil
	}

	f, ok := m.nodes[from.UID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	t, ok := m.nodes[to.UID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	// make a copy of link attributes and metadata
	attrs := attrs.New()
	metadata := metadata.New()

	if opts.Attrs != nil {
		for _, k := range opts.Attrs.Keys() {
			//fmt.Println("Setting attribute", k, "to value", opts.Attrs.Get(k))
			attrs.Set(k, opts.Attrs.Get(k))
		}
	}

	if opts.Metadata != nil {
		for _, k := range opts.Metadata.Keys() {
			metadata.Set(k, opts.Metadata.Get(k))
		}
	}

	eopts := []entity.Option{
		entity.Attrs(attrs),
		entity.Metadata(metadata),
		entity.Weight(opts.Weight),
		entity.Relation(opts.Relation),
	}

	uid := uuid.New().String()
	edge := NewEdge(uid, f, t, eopts...)

	m.SetWeightedEdge(edge)

	m.edges[uid] = edge

	return edge.Edge, nil
}

// SubGraph returns the subgraph of the node up to given depth or returns error
func (m *Memory) SubGraph(n store.Node, depth int) (store.Graph, error) {
	rootNode, ok := m.nodes[n.UID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	s := &Memory{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(1.0, 1.0),
		id:                      "subgraph" + m.id,
		nodes:                   make(map[string]*Node),
		edges:                   make(map[string]*Edge),
	}

	var sgErr error
	// k2g maps kraph node IDs to subgraph g nodes
	k2g := make(map[int64]store.Node)

	visit := func(n graph.Node) {
		vnode := n.(*Node)

		// create a deep copy of the Kraph node
		nodeAttrs := attrs.New()
		nodeMetadata := metadata.New()

		for _, k := range vnode.Attrs().Keys() {
			nodeAttrs.Set(k, vnode.Attrs().Get(k))
		}

		for _, k := range vnode.Metadata().Keys() {
			nodeMetadata.Set(k, vnode.Metadata().Get(k))
		}

		obj := vnode.Metadata().Get("object").(api.Object)
		opts := store.AddOptions{
			Attrs:    nodeAttrs,
			Metadata: nodeMetadata,
		}
		storeNode, err := s.Add(obj, opts)
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
				if _, err := s.Edge(storeNode.UID(), to.UID()); goerr.Is(err, errors.ErrEdgeNotExist) {
					// get the original edge from the memory store
					medge, err := m.Edge(vnode.UID(), kraphPeer.UID())
					if goerr.Is(err, errors.ErrEdgeNotExist) {
						sgErr = err
						return
					}

					attrs := attrs.New()
					metadata := metadata.New()

					for _, k := range medge.Attrs().Keys() {
						attrs.Set(k, medge.Attrs().Get(k))
					}

					for _, k := range medge.Metadata().Keys() {
						metadata.Set(k, medge.Metadata().Get(k))
					}

					opts := store.LinkOptions{
						Weight:   medge.Weight(),
						Attrs:    attrs,
						Metadata: metadata,
					}

					if _, err = s.Link(storeNode, to, opts); err != nil {
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

// DOT returns the GrapViz dot representation of kraph.
func (m *Memory) DOT() (string, error) {
	b, err := dot.Marshal(m.WeightedUndirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode into DOT: %w", err)
	}

	return string(b), nil
}
