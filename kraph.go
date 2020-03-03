package kraph

import (
	"errors"
	"fmt"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	// ErrNotImplemented is returned by functions whose functionality has not been implemented yet
	ErrNotImplemented = errors.New("not implemented")
	// ErrUnknownObject is returned when a given object is not recognised
	ErrUnknownObject = errors.New("unknown object")
	// ErrUnknownTop is returns when a given topology is not recognised
	ErrUnknownTop = errors.New("unknown topology")
)

// Kraph is a graph of Kubernetes resources
type Kraph struct {
	*simple.WeightedUndirectedGraph
	// client discovers and maps APIs
	client api.Client
	// options
	opts Options
	// Global DOT attributes
	GraphAttrs Attrs
	NodeAttrs  Attrs
	EdgeAttrs  Attrs
}

// New creates new Kraph with given options and returns it
// It never returns error, but it might in the future
func New(client api.Client, opts ...Option) (*Kraph, error) {
	kraphOpts := Options{}
	for _, apply := range opts {
		apply(&kraphOpts)
	}

	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		client:                  client,
		opts:                    kraphOpts,
		GraphAttrs:              make(Attrs),
		NodeAttrs:               make(Attrs),
		EdgeAttrs:               make(Attrs),
	}, nil
}

// Options returns kraph options
func (k *Kraph) Options() Options {
	return k.opts
}

// NewNode creates new kraph node and returns it as a gonum graph node
func (k *Kraph) NewNode(name string, opts ...NodeOption) *Node {
	nodeOpts := newNodeOptions(opts...)

	n := &Node{
		Attrs:    nodeOpts.Attrs,
		id:       k.WeightedUndirectedGraph.NewNode().ID(),
		name:     name,
		metadata: nodeOpts.Metadata,
	}

	for _, attr := range nodeOpts.Attrs.Attributes() {
		n.SetAttribute(attr.Key, attr.Value)
	}

	k.AddNode(n)

	return n
}

// NewEdge adds a new edge from source node to destination node to the graph
// or returns an existing edge if it already exists in the graph
// It will panic if the IDs of the from and to are equal
func (k *Kraph) NewEdge(from, to graph.Node, opts ...EdgeOption) *Edge {
	if e := k.Edge(from.ID(), to.ID()); e != nil {
		return e.(*Edge)
	}

	edgeOpts := newEdgeOptions(opts...)

	e := &Edge{
		Attrs:    edgeOpts.Attrs,
		from:     from.(*Node),
		to:       to.(*Node),
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

// DOTAttributers returns the global DOT kraph attributers
func (k *Kraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return k.GraphAttrs, k.NodeAttrs, k.EdgeAttrs
}

// DOT returns the GrapViz dot representation of kraph
func (k *Kraph) DOT(g graph.Graph) (string, error) {
	b, err := dot.Marshal(g, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode kraph into dot: %v", err)
	}

	return string(b), nil
}

func (k *Kraph) linkObjects(obj api.Object, neighbs []api.Object) {
	from := k.NewNode(obj.Name())

	for _, o := range neighbs {
		to := k.NewNode(o.Name())
		if e := k.Edge(from.ID(), to.ID()); e == nil {
			k.NewEdge(from, to)
		}
	}
}

func (k *Kraph) buildGraph(top api.Top) (graph.Graph, error) {
	switch r := top.Raw().(type) {
	// TODO: make this less hacky
	// One of the options is getting all objects
	// and then by iterating over them querying
	// the topology one by one when building the graph
	case map[string]map[string]map[string]api.Object:
		for _, kinds := range r {
			for _, names := range kinds {
				for _, obj := range names {
					raw := obj.Raw().(unstructured.Unstructured)
					var neighbs []api.Object
					for _, owner := range raw.GetOwnerReferences() {
						queryOpts := []query.Option{query.Kind(owner.Kind), query.Name(owner.Name)}
						objs, err := top.Get(queryOpts...)
						if err != nil {
							return nil, err
						}
						neighbs = append(neighbs, objs...)
					}
					k.linkObjects(obj, neighbs)
				}
			}
		}
	default:
		return nil, ErrUnknownTop
	}

	return nil, ErrNotImplemented
}

// Build builds resource graph and returns it
func (k *Kraph) Build() (graph.Graph, error) {
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

// Query allows to query for a kraph node
func (k *Kraph) Query(ns, kind, name string) (*Node, error) {
	return nil, ErrNotImplemented
}

// SubGraph returns a subgraph of given node up to given depth
func (k *Kraph) SubGraph(n *Node, depth int) (graph.Graph, error) {
	return nil, ErrNotImplemented
}

// GetNodesWithAttr returns a slice of nodes with the given attribute set
// If it does not find any matching nodes it returns empty slice.
// It returns error if attribute key is empty string,
func (k *Kraph) GetNodesWithAttr(attr encoding.Attribute) ([]*Node, error) {
	var nodes []*Node

	found := false
	for _, node := range graph.NodesOf(k.Nodes()) {
		n := node.(*Node)
		if val := n.Get(attr.Key); val != "" {
			// attribute key exists; check its value
			switch attr.Value {
			case "*":
				found = true
			case val:
				found = true
			default:
				// continue
			}
		}
		if found {
			nodes = append(nodes, n)
			found = false
		}
	}

	return nodes, nil
}

// GetEdgesWithAttr returns a slice of Edges with the given attribute set
// If it does not find any matching edges it returns empty slice.
// It returns error if attribute key is empty string,
func (k *Kraph) GetEdgesWithAttr(attr encoding.Attribute) ([]*Edge, error) {
	var edges []*Edge

	if attr.Key == "" {
		return nil, ErrAttrKeyInvalid
	}

	found := false
	for _, edge := range graph.EdgesOf(k.Edges()) {
		e := edge.(*Edge)
		if val := e.Get(attr.Key); val != "" {
			// attribute key exists; check its value
			switch attr.Value {
			case "*":
				found = true
			case val:
				found = true
			default:
				// continue
			}
		}
		if found {
			edges = append(edges, e)
			found = false
		}
	}

	return edges, nil
}
