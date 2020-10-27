package memory

import (
	"fmt"
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/attrs"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/metadata"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"github.com/milosgajdos/kraph/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/traverse"
)

// Memory is in-memory graph store
type Memory struct {
	// g is in-memory graph
	g *multi.WeightedUndirectedGraph
	// id is the store id
	id string
	// nodes maps api.Objects to graph Nodes
	nodes map[string]*Node
	// lines maps api.Object links to graph Edges
	lines map[string]*Line
	// options are store options
	opts store.Options
}

// NewStore creates new in-memory store and returns it
func NewStore(id string, opts store.Options) (*Memory, error) {
	return &Memory{
		g:     multi.NewWeightedUndirectedGraph(),
		id:    id,
		nodes: make(map[string]*Node),
		lines: make(map[string]*Line),
		opts:  opts,
	}, nil
}

// ID returns store ID
func (m Memory) ID() string {
	return m.id
}

// Options returns store options
func (m Memory) Options() store.Options {
	return m.opts
}

// Add adds obj to the store and returns it
func (m *Memory) Add(obj api.Object, opts store.AddOptions) (store.Entity, error) {
	uid := obj.UID().String()

	if node, ok := m.nodes[uid]; ok {
		return node, nil
	}

	var entOpts []entity.Option

	attrs := attrs.New()
	if opts.Attrs != nil {
		entOpts = append(entOpts, entity.Attrs(attrs))
	}

	metadata := metadata.New()
	if opts.Metadata != nil {
		entOpts = append(entOpts, entity.Metadata(metadata))
	}

	if obj.Resource() == nil {
		return nil, errors.ErrMissingResource
	}

	dotid := strings.Join([]string{
		obj.Resource().Version(),
		obj.Namespace(),
		obj.Resource().Kind(),
		obj.Name()}, "/")
	attrs.Set("name", dotid)

	n := m.g.NewNode()

	node := NewNode(n.ID(), uid, dotid, entOpts...)

	node.Metadata().Set("object", obj)

	m.g.AddNode(node)

	m.nodes[uid] = node

	return node, nil
}

// Delete deletes entity e from the memory store
func (m *Memory) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case store.Edge:
		l, ok := m.lines[e.UID()]
		if !ok {
			return fmt.Errorf("Edge Delete %s: %w", e.UID(), errors.ErrEdgeNotFound)
		}

		from, ok := m.nodes[v.From().UID()]
		if !ok {
			return fmt.Errorf("Edge Delete %s: %w", v.From().UID(), errors.ErrNodeNotFound)
		}

		to, ok := m.nodes[v.To().UID()]
		if !ok {
			return fmt.Errorf("Edge Delete %s: %w", v.To().UID(), errors.ErrNodeNotFound)
		}

		m.g.RemoveLine(from.ID(), to.ID(), l.ID())
		delete(m.lines, e.UID())
	case store.Node:
		node, ok := m.nodes[v.UID()]
		if !ok {
			return fmt.Errorf("Node Delete %s: %w", v.UID(), errors.ErrNodeNotFound)
		}

		m.g.RemoveNode(node.ID())
		delete(m.nodes, v.UID())
	default:
		return errors.ErrUnknownEntity
	}

	return nil
}

// QueryNode returns all the nodes that match given query.
func (m *Memory) QueryNode(q *query.Query) ([]*Node, error) {
	match := q.Matcher()

	if quid := match.UID(); quid != nil {
		if uid, ok := quid.Value().(uuid.UID); ok && len(uid.String()) > 0 {
			if n, ok := m.nodes[uid.String()]; ok {
				return []*Node{n}, nil
			}
		}
	}

	var results []*Node

	visit := func(n graph.Node) {
		node := n.(*Node)
		nodeObj := node.Metadata().Get("object").(api.Object)

		if match.NamespaceVal(nodeObj.Namespace()) {
			if match.KindVal(nodeObj.Resource().Kind()) {
				if match.NameVal(nodeObj.Name()) {
					if !match.AttrsVal(node.Attrs()) {
						return
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

					dotid := strings.Join([]string{
						nodeObj.Resource().Version(),
						nodeObj.Namespace(),
						nodeObj.Resource().Kind(),
						nodeObj.Name()}, "/")
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

	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	dfs.WalkAll(m.g, nil, nil, func(graph.Node) {})

	return results, nil
}

// QueryEdge returns all the edges that match given query
func (m *Memory) QueryLine(q *query.Query) ([]*Line, error) {
	match := q.Matcher()

	var results []*Line

	if quid := match.UID(); quid != nil {
		if uid, ok := quid.Value().(string); ok && len(uid) > 0 {
			if l, ok := m.lines[uid]; ok {
				return []*Line{l}, nil
			}
		}
	}

	trav := func(e graph.Edge) bool {
		from := e.From().(*Node)
		to := e.To().(*Node)

		if lines := m.g.WeightedLines(from.ID(), to.ID()); lines != nil {
			for lines.Next() {
				wl := lines.WeightedLine()
				we := wl.(*Line).Edge
				if match.WeightVal(we.Weight()) {
					if !match.AttrsVal(we.Attrs()) {
						continue
					}

					attrs := attrs.New()
					metadata := metadata.New()

					for _, k := range we.Attrs().Keys() {
						attrs.Set(k, we.Attrs().Get(k))
					}

					for _, k := range we.Metadata().Keys() {
						metadata.Set(k, we.Metadata().Get(k))
					}

					opts := []entity.Option{
						entity.Attrs(attrs),
						entity.Metadata(metadata),
						entity.Weight(we.Weight()),
					}

					ent := NewLine(wl.ID(), we.UID(), we.UID(), from, to, opts...)

					results = append(results, ent)
				}
			}
		}

		return true
	}

	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	dfs.WalkAll(m.g, nil, nil, func(graph.Node) {})

	return results, nil
}

// Query queries the in-memory graph and returns the matched results.
func (m *Memory) Query(q *query.Query) ([]store.Entity, error) {
	var e query.Entity

	if m := q.Matcher().Entity(); m != nil {
		var ok bool
		e, ok = m.Value().(query.Entity)
		if !ok {
			return nil, errors.ErrInvalidEntity
		}
	}

	var entities []store.Entity

	switch e {
	case query.Node:
		nodes, err := m.QueryNode(q)
		if err != nil {
			return nil, fmt.Errorf("Node query: %w", err)
		}
		for _, node := range nodes {
			entities = append(entities, node.Node)
		}
	case query.Edge:
		edges, err := m.QueryLine(q)
		if err != nil {
			return nil, fmt.Errorf("Edge query: %w", err)
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
	graphNodes := graph.NodesOf(m.g.Nodes())

	nodes := make([]store.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Edges returns all the edges (lines) from u to v
// if such edges exists and nil otherwise
func (m *Memory) Edges(uid, vid string) ([]store.Edge, error) {
	from, ok := m.nodes[uid]
	if !ok {
		return nil, fmt.Errorf("Edges %s: %w", uid, errors.ErrNodeNotFound)
	}

	to, ok := m.nodes[vid]
	if !ok {
		return nil, fmt.Errorf("Edges %s: %w", vid, errors.ErrNodeNotFound)
	}

	var edges []store.Edge

	if lines := m.g.WeightedLines(from.ID(), to.ID()); lines != nil {
		for lines.Next() {
			wl := lines.WeightedLine()
			we := wl.(*Line).Edge
			edges = append(edges, we)
		}

		return edges, nil
	}

	return nil, errors.ErrEdgeNotExist
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It returns error if either of the nodes does not exist in the graph.
func (m *Memory) Link(from store.Node, to store.Node, opts store.LinkOptions) (store.Edge, error) {
	f, ok := m.nodes[from.UID()]
	if !ok {
		return nil, fmt.Errorf("Link %s: %w", from.UID(), errors.ErrNodeNotFound)
	}

	t, ok := m.nodes[to.UID()]
	if !ok {
		return nil, fmt.Errorf("Link %s: %w", to.UID(), errors.ErrNodeNotFound)
	}

	if we := m.g.WeightedEdge(f.ID(), t.ID()); we != nil {
		if !opts.Line {
			me := we.(multi.WeightedEdge)
			if me.Next() {
				wl := me.WeightedLine()
				return wl.(*Line).Edge, nil
			}
		}
	}

	var entOpts []entity.Option

	attrs := attrs.New()
	if opts.Attrs != nil {
		for _, k := range opts.Attrs.Keys() {
			attrs.Set(k, opts.Attrs.Get(k))
		}
	}
	entOpts = append(entOpts, entity.Attrs(attrs))

	metadata := metadata.New()
	if opts.Metadata != nil {
		for _, k := range opts.Metadata.Keys() {
			metadata.Set(k, metadata.Get(k))
		}
	}
	entOpts = append(entOpts, entity.Metadata(metadata))

	if len(opts.Relation) > 0 {
		entOpts = append(entOpts, entity.Relation(opts.Relation))
	}

	w := opts.Weight
	if opts.Weight < 0 {
		w = store.DefaultWeight
	}
	entOpts = append(entOpts, entity.Weight(w))

	wl := m.g.NewWeightedLine(f, t, w)

	uid := uuid.New().String()

	line := NewLine(wl.ID(), uid, uid, f, t, entOpts...)

	m.g.SetWeightedLine(line)

	m.lines[uid] = line

	return line.Edge, nil
}

// SubGraph returns the subgraph of the node up to given depth
func (m *Memory) SubGraph(n store.Node, depth int) (store.Graph, error) {
	rootNode, ok := m.nodes[n.UID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	s := &Memory{
		g:     multi.NewWeightedUndirectedGraph(),
		id:    "sub-" + m.id,
		nodes: make(map[string]*Node),
		lines: make(map[string]*Line),
	}

	var sgErr error

	subnodes := make(map[int64]store.Node)

	visit := func(n graph.Node) {
		vnode := n.(*Node)

		nodeAttrs := attrs.New()
		for _, k := range vnode.Attrs().Keys() {
			nodeAttrs.Set(k, vnode.Attrs().Get(k))
		}

		nodeMetadata := metadata.New()
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

		subnodes[n.ID()] = storeNode
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	_ = bfs.Walk(m.g, rootNode, func(n graph.Node, d int) bool {
		if d == depth {
			return true
		}
		return false
	})

	if sgErr != nil {
		return nil, sgErr
	}

	for id, node := range subnodes {
		nodes := s.g.From(id)
		for nodes.Next() {
			pnode := nodes.Node()
			peer := pnode.(*Node)
			if to, ok := subnodes[peer.ID()]; ok {
				if lines := m.g.WeightedLines(id, pnode.ID()); lines != nil {
					for lines.Next() {
						wl := lines.WeightedLine()
						we := wl.(*Line).Edge
						attrs := attrs.New()
						for _, k := range we.Attrs().Keys() {
							attrs.Set(k, we.Attrs().Get(k))
						}

						metadata := metadata.New()
						for _, k := range we.Metadata().Keys() {
							metadata.Set(k, we.Metadata().Get(k))
						}

						opts := store.LinkOptions{
							Line:     true,
							Weight:   we.Weight(),
							Attrs:    attrs,
							Metadata: metadata,
						}

						if _, err := s.Link(node, to, opts); err != nil {
							return nil, fmt.Errorf("Subgraph: %v", err)
						}
					}
				}
			}
		}
	}

	return s, nil
}

// DOTID returns the store DOT ID.
func (m *Memory) DOTID() string {
	return m.id
}

// DOTAttributers implements encoding.Attributer
func (m *Memory) DOTAttributers() (graph, node, edge encoding.Attributer) {
	graph = attrs.New()
	if m.opts.DOTOptions.GraphAttrs != nil {
		graph = m.opts.DOTOptions.GraphAttrs
	}

	node = attrs.New()
	if m.opts.DOTOptions.NodeAttrs != nil {
		node = m.opts.DOTOptions.NodeAttrs
	}

	edge = attrs.New()
	if m.opts.DOTOptions.EdgeAttrs != nil {
		edge = m.opts.DOTOptions.EdgeAttrs
	}

	return graph, node, edge
}

// DOT returns the GrapViz dot representation of kraph.
func (m *Memory) DOT() (string, error) {
	b, err := dot.MarshalMulti(m.g, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode into DOT: %w", err)
	}

	return string(b), nil
}
