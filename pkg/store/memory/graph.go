package memory

import (
	"fmt"
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/errors"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/entity"
	"github.com/milosgajdos/kraph/pkg/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/traverse"
)

// Graph is in-memory graph
type Graph struct {
	// g is in-memory graph
	*multi.WeightedUndirectedGraph
	// id is ID of the graph
	id string
	// nodes maps api.Objects to graph Nodes
	nodes map[string]*Node
	// lines maps api.Object links to graph Edges
	lines map[string]*Line
	// opts are graph options
	opts store.GraphOptions
}

// NewGraph creates a new graph and returns it.
func NewGraph(id string, opts store.GraphOptions) *Graph {
	return &Graph{
		WeightedUndirectedGraph: multi.NewWeightedUndirectedGraph(),
		id:                      id,
		nodes:                   make(map[string]*Node),
		lines:                   make(map[string]*Line),
		opts:                    opts,
	}
}

// NewNode adds a node to graph.
func (g *Graph) NewNode(obj api.Object, opts store.AddOptions) (*Node, error) {
	uid := obj.UID().String()

	if node, ok := g.nodes[uid]; ok {
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

	n := g.WeightedUndirectedGraph.NewNode()

	node := NewNode(n.ID(), uid, dotid, entOpts...)

	node.Metadata().Set("object", obj)

	g.nodes[uid] = node

	g.WeightedUndirectedGraph.AddNode(node)

	return node, nil
}

// RemoveNode removes the node with given id from graph.
func (g *Graph) RemoveNode(id string) error {
	node, ok := g.nodes[id]
	if !ok {
		return fmt.Errorf("Node Delete %s: %w", id, errors.ErrNodeNotFound)
	}

	g.WeightedUndirectedGraph.RemoveNode(node.ID())
	delete(g.nodes, id)

	return nil
}

// RemoveLine removes the line between two nodes.
func (g *Graph) RemoveLine(from, to, lid string) error {
	l, ok := g.lines[lid]
	if !ok {
		return fmt.Errorf("Edge Delete %s: %w", lid, errors.ErrEdgeNotFound)
	}

	f, ok := g.nodes[from]
	if !ok {
		return fmt.Errorf("Edge Delete %s: %w", from, errors.ErrNodeNotFound)
	}

	t, ok := g.nodes[to]
	if !ok {
		return fmt.Errorf("Edge Delete %s: %w", to, errors.ErrNodeNotFound)
	}

	g.WeightedUndirectedGraph.RemoveLine(f.ID(), t.ID(), l.ID())
	delete(g.lines, lid)

	return nil
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (g *Graph) Node(id string) (store.Node, error) {
	if node, ok := g.nodes[id]; ok {
		return node, nil
	}

	return nil, errors.ErrNodeNotFound
}

// Nodes returns all the nodes in the graph.
func (g *Graph) Nodes() ([]store.Node, error) {
	graphNodes := graph.NodesOf(g.WeightedUndirectedGraph.Nodes())

	nodes := make([]store.Node, len(graphNodes))

	for i, n := range graphNodes {
		nodes[i] = n.(*Node)
	}

	return nodes, nil
}

// Edges returns all the edges (lines) from u to v
// if such edges exists and nil otherwise
func (g *Graph) Edges(uid, vid string) ([]store.Edge, error) {
	from, ok := g.nodes[uid]
	if !ok {
		return nil, fmt.Errorf("Edges %s: %w", uid, errors.ErrNodeNotFound)
	}

	to, ok := g.nodes[vid]
	if !ok {
		return nil, fmt.Errorf("Edges %s: %w", vid, errors.ErrNodeNotFound)
	}

	var edges []store.Edge

	if lines := g.WeightedLines(from.ID(), to.ID()); lines != nil {
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
func (g *Graph) Link(from store.Node, to store.Node, opts store.LinkOptions) (store.Edge, error) {
	f, ok := g.nodes[from.UID()]
	if !ok {
		return nil, fmt.Errorf("Link %s: %w", from.UID(), errors.ErrNodeNotFound)
	}

	t, ok := g.nodes[to.UID()]
	if !ok {
		return nil, fmt.Errorf("Link %s: %w", to.UID(), errors.ErrNodeNotFound)
	}

	if we := g.WeightedEdge(f.ID(), t.ID()); we != nil {
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

	wl := g.NewWeightedLine(f, t, w)

	uid := uuid.New().String()

	line := NewLine(wl.ID(), uid, uid, f, t, entOpts...)

	g.SetWeightedLine(line)

	g.lines[uid] = line

	return line.Edge, nil
}

// SubGraph returns the subgraph of the node up to given depth
func (g *Graph) SubGraph(n store.Node, depth int) (store.Graph, error) {
	rootNode, ok := g.nodes[n.UID()]
	if !ok {
		return nil, errors.ErrNodeNotFound
	}

	sub := NewGraph(g.id+"subgraph", g.opts)

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

		if _, err := sub.NewNode(obj, opts); err != nil {
			sgErr = err
			return
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	_ = bfs.Walk(g, rootNode, func(n graph.Node, d int) bool {
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
				if lines := g.WeightedLines(id, pnode.ID()); lines != nil {
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

						if _, err := sub.Link(node, to, opts); err != nil {
							return nil, fmt.Errorf("Subgraph: %v", err)
						}
					}
				}
			}
		}
	}

	return sub, nil
}

// DOTID returns the store DOT ID.
func (g Graph) DOTID() string {
	return g.id
}

// DOTAttributers implements encoding.Attributer
func (g *Graph) DOTAttributers() (graph, node, edge encoding.Attributer) {
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
func (g *Graph) DOT() (string, error) {
	b, err := dot.MarshalMulti(g.WeightedUndirectedGraph, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode into DOT: %w", err)
	}

	return string(b), nil
}
