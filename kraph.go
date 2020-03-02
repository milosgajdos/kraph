package kraph

import (
	"errors"
	"fmt"

	"github.com/milosgajdos/kraph/api"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	// ErrNotImplemented is returned by functions whose functionality has not been implemented yet
	ErrNotImplemented = errors.New("not implemented")
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
func (k *Kraph) DOT() (string, error) {
	b, err := dot.Marshal(k, "", "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode kraph into dot: %v", err)
	}

	return string(b), nil
}

// Build builds resource graph
func (k *Kraph) Build() error {
	api, err := k.client.Discover()
	if err != nil {
		return fmt.Errorf("failed discovering API: %w", err)
	}

	if err := k.client.Map(api); err != nil {
		return fmt.Errorf("failed mapping API: %w", err)
	}

	return nil
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

//// addEdge creates an between from and to nodes with given weight
//func (k *Kraph) addEdge(from, to graph.Node, w float64) {
//	var a Attrs = map[string]string{
//		"weight": fmt.Sprintf("%f", w),
//	}
//
//	k.NewEdge(from, to, EdgeAttrs(a), Weight(w))
//}
//
//// addNode adds a new node to the graph and returns it
//func (k *Kraph) addNode(res interface{}) *Node {
//	var name, kind string
//	var uid types.UID
//	var m Metadata
//
//	switch r := res.(type) {
//	case unstructured.Unstructured:
//		name = resName(r.GetKind(), r.GetName())
//		kind = strings.ToLower(r.GetKind())
//		uid = resUID(r)
//		if ns := r.GetNamespace(); ns != "" {
//			if k.nsMap[ns] == nil {
//				k.nsMap[ns] = make(map[string]types.UID)
//			}
//			k.nsMap[ns][name] = uid
//		}
//	case metav1.OwnerReference:
//		name = resName(r.Kind, r.Name)
//		kind = strings.ToLower(r.Kind)
//		uid = r.UID
//	}
//
//	var a Attrs = map[string]string{
//		"kind": kind,
//	}
//	m = map[string]interface{}{
//		"resource": res,
//	}
//	node := k.NewNode(name, NodeAttrs(a), NodeMetadata(m))
//
//	k.nodeMap[kind][uid] = node
//
//	return node
//}
//
//// linkNode links the node too all of its owners
//func (k *Kraph) linkNode(from graph.Node) {
//	knode := from.(*Node)
//	res := knode.metadata["resource"].(unstructured.Unstructured)
//
//	for _, owner := range res.GetOwnerReferences() {
//		kind := strings.ToLower(owner.Kind)
//		uid := owner.UID
//		if k.nodeMap[kind] == nil {
//			k.nodeMap[kind] = make(map[types.UID]*Node)
//		}
//		if to, ok := k.nodeMap[kind][uid]; ok {
//			if e := k.Edge(from.ID(), to.ID()); e == nil {
//				k.addEdge(from, to, 0.0)
//			}
//			continue
//		}
//
//		to := k.addNode(owner)
//
//		k.addEdge(from, to, 0.0)
//	}
//}
//
//// linkResource links API resource to all of its owners
//func (k *Kraph) linkResource(res unstructured.Unstructured) {
//	kind := strings.ToLower(res.GetKind())
//	if k.nodeMap[kind] == nil {
//		k.nodeMap[kind] = make(map[types.UID]*Node)
//	}
//
//	uid := resUID(res)
//
//	// we need to check for the existence of the node in nodeMap
//	// as the order of discovered API objects is arbitrary so some
//	// resources may have linked themselves into their Owner refs
//	// and added those Owner refs as "placeholder" nodes into the nodeMap
//	if n, ok := k.nodeMap[kind][uid]; ok {
//		n.metadata["resource"] = res
//		return
//	}
//
//	node := k.addNode(res)
//
//	k.linkNode(node)
//}
