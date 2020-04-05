package kraph

import (
	"fmt"
	"math/big"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

var (
	// DefaultWeight is the default edge weight
	DefaultWeight = 0.0
)

// Kraph is a graph of Kubernetes resources
type Kraph struct {
	*simple.WeightedUndirectedGraph
	// nodes maps api.Objects into their node.ID
	nodes map[string]int64
	// client discovers and maps APIs
	client api.Client
	// options
	opts Options
	// Global DOT attributes
	GraphAttrs Attrs
	NodeAttrs  Attrs
	EdgeAttrs  Attrs
}

// New creates new Kraph with given options and returns it.
// It never returns error at the moment, but it might in the future.
func New(client api.Client, opts ...Option) (*Kraph, error) {
	kraphOpts := Options{}
	for _, apply := range opts {
		apply(&kraphOpts)
	}

	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		nodes:                   make(map[string]int64),
		client:                  client,
		opts:                    kraphOpts,
		GraphAttrs:              make(Attrs),
		NodeAttrs:               make(Attrs),
		EdgeAttrs:               make(Attrs),
	}, nil
}

// Options returns kraph options.
func (k *Kraph) Options() Options {
	return k.opts
}

// NewNode creates new kraph node, adds it to its graph and returns it.
func (k *Kraph) NewNode(obj api.Object, opts ...NodeOption) *Node {
	if id, ok := k.nodes[obj.UID().String()]; ok {
		node := k.Node(id)
		return node.(*Node)
	}

	nodeOpts := newNodeOptions(opts...)

	n := &Node{
		Attrs:    nodeOpts.Attrs,
		id:       k.WeightedUndirectedGraph.NewNode().ID(),
		name:     obj.Kind() + "-" + obj.Name(),
		metadata: nodeOpts.Metadata,
	}

	for _, attr := range nodeOpts.Attrs.Attributes() {
		n.SetAttribute(attr.Key, attr.Value)
	}

	n.metadata["object"] = obj

	k.AddNode(n)

	k.nodes[obj.UID().String()] = n.ID()

	return n
}

// NewEdge adds a new edge between from and to nodes to kraph
// or returns an existing edge if it already exists in the graph.
// It will panic if the IDs of the from and to nodes are the same.
//func (k *Kraph) NewEdge(from, to graph.Node, opts ...EdgeOption) *Edge {
func (k *Kraph) NewEdge(from, to *Node, opts ...EdgeOption) *Edge {
	if e := k.Edge(from.ID(), to.ID()); e != nil {
		return e.(*Edge)
	}

	edgeOpts := newEdgeOptions(opts...)

	e := &Edge{
		Attrs:    edgeOpts.Attrs,
		from:     from,
		to:       to,
		weight:   edgeOpts.Weight,
		metadata: edgeOpts.Metadata,
	}

	for _, attr := range edgeOpts.Attrs.Attributes() {
		e.SetAttribute(attr.Key, attr.Value)
	}

	k.SetWeightedEdge(e)

	return e
}

// DOTID returns the graph's DOT ID.
func (k *Kraph) DOTID() string {
	return "kraph"
}

// DOTAttributers returns the global DOT kraph attributers.
func (k *Kraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return k.GraphAttrs, k.NodeAttrs, k.EdgeAttrs
}

// DOT returns the GrapViz dot representation of kraph.
func (k *Kraph) DOT() (string, error) {
	b, err := dot.Marshal(k, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode kraph into dot: %v", err)
	}

	return string(b), nil
}

// linkObject links obj to all of its neighbours and sets their relation to rel.
func (k *Kraph) linkObjects(obj api.Object, rel api.Relation, neighbs []api.Object) {
	from := k.NewNode(obj)

	for _, o := range neighbs {
		to := k.NewNode(o)
		if e := k.Edge(from.ID(), to.ID()); e == nil {
			attrs := make(Attrs)
			attrs["relation"] = rel.String()
			e = k.NewEdge(from, to, EdgeAttrs(attrs))
		}
	}
}

// buildGraph builds a graph from given topology and returns it.
func (k *Kraph) buildGraph(top api.Top) (graph.Graph, error) {
	for _, object := range top.Objects() {
		if len(object.Links()) == 0 {
			k.NewNode(object)
			continue
		}
		for _, link := range object.Links() {
			query := []query.Option{
				query.UID(link.To().String()),
			}
			objs, err := top.Get(query...)
			if err != nil {
				return nil, err
			}
			k.linkObjects(object, link.Relation(), objs)
		}
	}

	return k.WeightedUndirectedGraph, nil
}

// Build builds resource graph and returns it.
func (k *Kraph) Build() (graph.Graph, error) {
	// TODO: reset the graph before building
	// This will allow to run Build multiple times
	// each time building the graph from scratch
	api, err := k.client.Discover()
	if err != nil {
		return nil, fmt.Errorf("failed discovering API: %w", err)
	}

	top, err := k.client.Map(api)
	if err != nil {
		return nil, fmt.Errorf("failed mapping API: %w", err)
	}

	return k.buildGraph(top)
}

// QueryNode returns all the nodes that match given query.
func (k *Kraph) QueryNode(opts ...query.Option) ([]*Node, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	if len(query.UID) > 0 {
		if id, ok := k.nodes[query.UID]; ok {
			node := k.Node(id)

			return []*Node{node.(*Node)}, nil
		}
	}

	var results []*Node

	visit := func(n graph.Node) {
		node := n.(*Node)
		nodeObj := node.metadata["object"].(api.Object)

		if len(query.Namespace) == 0 || query.Namespace == nodeObj.Namespace() {
			if len(query.Kind) == 0 || query.Kind == nodeObj.Kind() {
				if len(query.Name) == 0 || query.Name == nodeObj.Name() {
					if len(query.Attrs) > 0 {
						for k, v := range query.Attrs {
							if node.GetAttribute(k) != v {
								return
							}
						}
					}

					// create a deep copy of the matched node
					attrs := make(Attrs)
					metadata := make(Metadata)

					for k, v := range node.Attrs {
						attrs.SetAttribute(k, v)
					}

					for k, v := range node.metadata {
						metadata[k] = v
					}

					qNode := &Node{
						Attrs:    attrs,
						id:       node.ID(),
						name:     node.name,
						metadata: metadata,
					}
					results = append(results, qNode)
				}
			}
		}
	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(k, nil, nil, func(graph.Node) {})

	return results, nil
}

// QueryEdge returns all the edges that match given query
func (k *Kraph) QueryEdge(opts ...query.Option) ([]*Edge, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var results []*Edge

	traversed := make(map[int64]map[int64]bool)

	trav := func(e graph.Edge) bool {
		edge := e.(*Edge)

		if traversed[edge.from.ID()] == nil {
			traversed[edge.from.ID()] = make(map[int64]bool)
		}

		if traversed[edge.to.ID()] == nil {
			traversed[edge.to.ID()] = make(map[int64]bool)
		}

		if traversed[edge.from.ID()][edge.to.ID()] || traversed[edge.to.ID()][edge.from.ID()] {
			return false
		}

		traversed[edge.from.ID()][edge.to.ID()] = true
		traversed[edge.to.ID()][edge.from.ID()] = true

		if big.NewFloat(query.Weight).Cmp(big.NewFloat(edge.weight)) == 0 {
			if len(query.Attrs) > 0 {
				for k, v := range query.Attrs {
					if edge.GetAttribute(k) != v {
						return false
					}
				}
			}

			// create a deep copy of the matched edge
			attrs := make(Attrs)
			metadata := make(Metadata)

			for k, v := range edge.Attrs {
				attrs.SetAttribute(k, v)
			}

			for k, v := range edge.metadata {
				metadata[k] = v
			}

			qEdge := &Edge{
				Attrs:    attrs,
				from:     edge.from,
				to:       edge.to,
				weight:   edge.weight,
				metadata: edge.metadata,
			}

			results = append(results, qEdge)
		}

		return true

	}

	// let's go with DFS as it's more memory efficient
	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	// traverse the whole graph and collect all nodes matching the query
	dfs.WalkAll(k, nil, nil, func(graph.Node) {})

	return results, nil
}

// SubGraph returns a subgraph of node n up to given depth.
// It performs a Breadth First Search (BFS) and creates a subgraph
// from the nodes traversed during the search.
func (k *Kraph) SubGraph(n *Node, depth int) (graph.Graph, error) {
	g := simple.NewWeightedUndirectedGraph(0.0, 0.0)

	// k2g maps kraph node IDs to subgraph g nodes
	k2g := make(map[int64]graph.Node)

	visit := func(n graph.Node) {
		node := n.(*Node)

		// create a deep copy of the Kraph node
		attrs := make(Attrs)
		metadata := make(Metadata)

		for k, v := range node.Attrs {
			attrs.SetAttribute(k, v)
		}

		for k, v := range node.metadata {
			metadata[k] = v
		}

		gNode := &Node{
			Attrs:    attrs,
			id:       g.NewNode().ID(),
			name:     node.name,
			metadata: metadata,
		}

		g.AddNode(gNode)
		k2g[n.ID()] = gNode

		// this is not very efficient
		// the idea here is we go through newly visited node
		// and check if any of its peer nodes from Kraph have
		// been visited (k2g map) and if yes, then wire them
		// to this newly created subgraph node if they
		// have not already been wired to this node
		nodes := k.From(n.ID())
		for nodes.Next() {
			kraphPeer := nodes.Node()
			if to, ok := k2g[kraphPeer.ID()]; ok {
				if e := g.Edge(gNode.ID(), to.ID()); e == nil {
					edge := k.Edge(n.ID(), kraphPeer.ID())
					kEdge := edge.(*Edge)

					attrs := make(Attrs)
					metadata := make(Metadata)

					for k, v := range node.Attrs {
						attrs.SetAttribute(k, v)
					}

					for k, v := range node.metadata {
						metadata[k] = v
					}

					e := &Edge{
						Attrs:    attrs,
						from:     gNode,
						to:       to.(*Node),
						weight:   kEdge.weight,
						metadata: metadata,
					}

					g.SetWeightedEdge(e)
				}
			}
		}
	}

	bfs := traverse.BreadthFirst{
		Visit: visit,
	}

	// keep traversing until you reach the requested depth
	_ = bfs.Walk(k, n, func(n graph.Node, d int) bool {
		if d == depth {
			return true
		}
		return false
	})

	return g, nil
}
