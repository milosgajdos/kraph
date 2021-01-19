package memory

import (
	"fmt"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/graph"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/uuid"
	gngraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

// WUG is in-memory graph
type WUG struct {
	// g is in-memory graph
	*simple.WeightedUndirectedGraph
	// id is ID of the graph
	id string
	// nodes maps graph nodes
	nodes map[string]graph.Node
	// opts are graph options
	opts graph.Options
}

// NewWUG creates a new weighted undirected graph and returns it.
func NewWUG(id string, opts graph.Options) (*WUG, error) {
	return &WUG{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(opts.Weight, opts.Weight),
		id:                      id,
		nodes:                   make(map[string]graph.Node),
		opts:                    opts,
	}, nil
}

// NewNode creates a new graph node and returns it.
// NOTE: this is a convenience method which both creates a new *Node and
// adds it to the underlying graph with a unique ID
// The alternative is calling g.NewNode() and NewNode(node.ID(), obj, opts...)
func (g *WUG) NewNode(obj api.Object, opts ...entity.Option) (graph.Node, error) {
	gnode := g.WeightedUndirectedGraph.NewNode()

	nodeOpts := entity.NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	node, err := NewNode(gnode.ID(), obj, opts...)
	if err != nil {
		return nil, err
	}

	if n, ok := g.nodes[node.UID()]; ok {
		return n, nil
	}

	return node, nil
}

// AddNode adds node to the graph.
func (g *WUG) AddNode(n graph.Node) error {
	if _, ok := g.nodes[n.UID()]; ok {
		return nil
	}

	gnode, ok := n.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	if node := g.WeightedUndirectedGraph.Node(gnode.ID()); node != nil {
		g.nodes[n.UID()] = n
		return nil
	}

	g.nodes[n.UID()] = n

	g.WeightedUndirectedGraph.AddNode(gnode)

	return nil
}

// RemoveNode removes the node with the given id from graph.
func (g *WUG) RemoveNode(id string) error {
	node, ok := g.nodes[id]
	if !ok {
		return nil
	}

	gnode, ok := node.(*Node)
	if !ok {
		return graph.ErrInvalidNode
	}

	g.WeightedUndirectedGraph.RemoveNode(gnode.ID())

	delete(g.nodes, id)

	return nil
}

// Node returns the node with the given ID if it exists
// in the graph, and error if it could not be retrieved.
func (g *WUG) Node(id string) (graph.Node, error) {
	if node, ok := g.nodes[id]; ok {
		return node, nil
	}

	return nil, graph.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (g *WUG) Nodes() ([]graph.Node, error) {
	graphNodes := gngraph.NodesOf(g.WeightedUndirectedGraph.Nodes())

	nodes := make([]graph.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Edge returns edge between the two nodes
func (g *WUG) Edge(uid, vid string) (graph.Edge, error) {
	from, ok := g.nodes[uid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", uid, graph.ErrNodeNotFound)
	}

	to, ok := g.nodes[vid]
	if !ok {
		return nil, fmt.Errorf("%s: %w", vid, graph.ErrNodeNotFound)
	}

	// NOTE: it's safe to typecast without checking as
	// AddNode() checks for the right type befoe adding the node to graph.
	if e := g.WeightedEdge(from.(*Node).ID(), to.(*Node).ID()); e != nil {
		return e.(*Edge), nil
	}

	return nil, graph.ErrEdgeNotExist
}

// Edges returns all the edges (lines) from u to v.
func (g *WUG) Edges() ([]graph.Edge, error) {
	wedges := g.WeightedUndirectedGraph.Edges()

	graphEdges := gngraph.EdgesOf(wedges)

	edges := make([]graph.Edge, len(graphEdges))

	for i, e := range graphEdges {
		edges[i] = e.(*Edge)
	}

	return edges, nil
}

// RemoveEdge removes the line between two nodes.
func (g *WUG) RemoveEdge(from, to string) error {
	f, ok := g.nodes[from]
	if !ok {
		return nil
	}

	t, ok := g.nodes[to]
	if !ok {
		return nil
	}

	// NOTE: it's safe to typecast without checking as
	// AddNode() checks for the right type befoe adding the node to graph.
	g.WeightedUndirectedGraph.RemoveEdge(f.(*Node).ID(), t.(*Node).ID())

	return nil
}

// Link creates a new edge between from and to and returns it or it returns the existing edge.
// It returns error if either of the nodes does not exist in the graph.
func (g *WUG) Link(from, to string, opts graph.LinkOptions) (graph.Edge, error) {
	e, err := g.Edge(from, to)
	if err != nil && err != graph.ErrEdgeNotExist {
		return nil, err
	}

	if e != nil {
		return e, nil
	}

	f, ok := g.nodes[from]
	if !ok {
		return nil, fmt.Errorf("node %s link error: %w", from, graph.ErrNodeNotFound)
	}

	t, ok := g.nodes[to]
	if !ok {
		return nil, fmt.Errorf("node %s link error: %w", to, graph.ErrNodeNotFound)
	}

	var entOpts []entity.Option

	attrs := attrs.NewCopyFrom(opts.Attrs)
	entOpts = append(entOpts, entity.Attrs(attrs))

	w := opts.Weight
	if opts.Weight < 0 {
		w = graph.DefaultWeight
	}

	edge, err := NewEdge(f.(*Node), t.(*Node), w, entOpts...)
	if err != nil {
		return nil, err
	}

	g.SetWeightedEdge(edge)

	return edge, nil
}

// SubGraph returns the subgraph of the node up to given depth.
func (g *WUG) SubGraph(uid string, depth int) (graph.Graph, error) {
	root, ok := g.nodes[uid]
	if !ok {
		return nil, graph.ErrNodeNotFound
	}

	sub, err := NewWUG(g.id+"subgraph", g.opts)
	if err != nil {
		return nil, err
	}

	var sgErr error

	subnodes := make(map[int64]graph.Node)

	visit := func(n gngraph.Node) {
		vnode := n.(*Node)

		if err := sub.AddNode(vnode); err != nil {
			sgErr = err
			return
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	_ = bfs.Walk(g.WeightedUndirectedGraph, root.(*Node), func(n gngraph.Node, d int) bool {
		if d == depth {
			return true
		}
		return false
	})

	if sgErr != nil {
		return nil, sgErr
	}

	for id, node := range subnodes {
		nodes := sub.From(id)
		for nodes.Next() {
			pnode := nodes.Node()
			peer := pnode.(*Node)
			if to, ok := subnodes[peer.ID()]; ok {
				if edges := g.WeightedEdges(); edges != nil {
					for edges.Next() {
						we := edges.WeightedEdge()
						e := we.(*Edge)

						a := attrs.NewCopyFrom(e.Attrs())

						opts := graph.LinkOptions{
							Attrs:  a,
							Weight: e.Weight(),
						}

						if _, err := sub.Link(node.UID(), to.UID(), opts); err != nil {
							return nil, fmt.Errorf("subgraph %s link error: %v", sub.id, err)
						}
					}
				}
			}
		}
	}

	return sub, nil
}

// QueryEdge returns all the edges that match given query
func (g WUG) QueryEdge(q *query.Query) ([]graph.Edge, error) {
	match := q.Matcher()

	traversed := make(map[string]bool)

	var results []graph.Edge

	trav := func(e gngraph.Edge) bool {
		edge := e.(*Edge)

		if traversed[edge.UID()] {
			return false
		}

		traversed[edge.UID()] = true

		if match.WeightVal(edge.Weight()) {
			if !match.AttrsVal(edge.Attrs()) {
				return false
			}
			results = append(results, edge)
		}

		return true
	}

	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	dfs.WalkAll(g.WeightedUndirectedGraph, nil, nil, func(gngraph.Node) {})

	return results, nil
}

// QueryNode returns all the nodes that match the given query.
func (g WUG) QueryNode(q *query.Query) ([]graph.Node, error) {
	match := q.Matcher()

	if quid := match.UID(); quid != nil {
		if uid, ok := quid.Value().(uuid.UID); ok && len(uid.String()) > 0 {
			if n, ok := g.nodes[uid.String()]; ok {
				return []graph.Node{n}, nil
			}
		}
	}

	var results []graph.Node

	visit := func(n gngraph.Node) {
		node := n.(*Node)

		o := node.Object()

		if match.NamespaceVal(o.Namespace()) {
			if match.KindVal(o.Resource().Kind()) {
				if match.NameVal(o.Name()) {
					if !match.AttrsVal(node.Attrs()) {
						return
					}

					results = append(results, node)
				}
			}
		}
	}

	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	dfs.WalkAll(g.WeightedUndirectedGraph, nil, nil, func(gngraph.Node) {})

	return results, nil
}

// Query queries the in-memory graph and returns the matched results.
func (g WUG) Query(q *query.Query) ([]graph.Entity, error) {
	var e query.Entity

	if m := q.Matcher().Entity(); m != nil {
		var ok bool
		e, ok = m.Value().(query.Entity)
		if !ok {
			return nil, graph.ErrUnknownEntity
		}
	}

	var entities []graph.Entity

	switch e {
	case query.Node:
		nodes, err := g.QueryNode(q)
		if err != nil {
			return nil, fmt.Errorf("node query: %w", err)
		}
		for _, node := range nodes {
			entities = append(entities, node)
		}
	case query.Edge:
		edges, err := g.QueryEdge(q)
		if err != nil {
			return nil, fmt.Errorf("edge query: %w", err)
		}
		for _, edge := range edges {
			entities = append(entities, edge)
		}
	default:
		return nil, graph.ErrUnknownEntity
	}

	return entities, nil
}

// DOTID returns the store DOT ID.
func (g WUG) DOTID() string {
	return g.id
}

// DOTAttributers implements encoding.Attributer.
func (g *WUG) DOTAttributers() (graph, node, edge encoding.Attributer) {
	graph = attrs.New()
	if g.opts.DOTOptions.GraphAttrs != nil {
		graph = g.opts.DOTOptions.GraphAttrs
	}

	node = attrs.New()
	if g.opts.DOTOptions.NodeAttrs != nil {
		node = g.opts.DOTOptions.NodeAttrs
	}

	edge = attrs.New()
	if g.opts.DOTOptions.EdgeAttrs != nil {
		edge = g.opts.DOTOptions.EdgeAttrs
	}

	return graph, node, edge
}

// DOT returns the GrapViz dot representation of kraph.
func (g *WUG) DOT() (string, error) {
	b, err := dot.Marshal(g.WeightedUndirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("DOT marshal error: %w", err)
	}

	return string(b), nil
}
