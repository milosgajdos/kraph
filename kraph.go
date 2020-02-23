package kraph

import (
	"context"
	"errors"
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrUnknownObject   = errors.New("unknown object")
	ErrInterruptedCall = errors.New("call interrupted")
)

// Node is graph node
type Node struct {
	graph.Node
	Attributes
	// Name names the node
	Name string
}

// DOTID returns the node's DOT ID.
func (n *Node) DOTID() string {
	return n.Name
}

// SetDOTID sets the node's DOT ID.
func (n *Node) SetDOTID(id string) {
	n.Name = id
}

// Edge is graph edge
type Edge struct {
	graph.WeightedEdge
	Attributes
}

// Kraph is a graph of Kubernetes resources
type Kraph struct {
	*simple.WeightedUndirectedGraph
	// nodeMap maps graph nodes to kubernetes IDs
	nodeMap map[types.UID]*Node
	// disc is kubernetes discovery client
	disc discovery.DiscoveryInterface
	// dyn is kubernetes dynamic client
	dyn dynamic.Interface
	// Global DOT attributes
	GraphAttrs Attributes
	NodeAttrs  Attributes
	EdgeAttrs  Attributes
}

// New creates new Kraph and returns it
//func New(client kubernetes.Interface) (*Kraph, error) {
func New(discover discovery.DiscoveryInterface, dynamic dynamic.Interface) (*Kraph, error) {
	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		nodeMap:                 make(map[types.UID]*Node),
		disc:                    discover,
		dyn:                     dynamic,
	}, nil
}

// NewNode creates new kraph node and returns it
func (k *Kraph) NewNode(name string) *Node {
	// TODO: should add the node to the graph for better UX
	// as if its not added in, it won't have unique ID
	return &Node{
		Node: k.WeightedUndirectedGraph.NewNode(),
		Name: name,
	}
}

// Nodes returns all kraph graph nodes
func (k *Kraph) Nodes() graph.Nodes {
	return k.WeightedUndirectedGraph.Nodes()
}

// NewEdge adds a new edge from source node to destination node to the graph
// or returns an existing edge if it already exists
// It will panic if the IDs of the from and to are equal
func (k *Kraph) NewEdge(from, to *Node, weight float64) *Edge {
	if e := k.Edge(from.ID(), to.ID()); e != nil {
		ke, ok := e.(*Edge)
		if !ok {
			return &Edge{
				WeightedEdge: e.(*simple.WeightedEdge),
			}
		}

		return ke
	}

	e := &Edge{
		WeightedEdge: k.WeightedUndirectedGraph.NewWeightedEdge(from, to, weight),
	}

	k.SetWeightedEdge(e)

	return e
}

// Edges returns all kraph graph edges
func (k *Kraph) Edges() graph.Edges {
	return k.WeightedUndirectedGraph.Edges()
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

// Build builds the kubernetes resource graph for a given namespace
// If the namespace is empty string all namespaces are assumed
func (k *Kraph) Build(ctx context.Context, namespace string) error {
	_, err := k.discoverAPI(ctx)
	if err != nil {
		return fmt.Errorf("failed discovering kubernetes API: %w", err)
	}

	return nil
}

// discoverAPI discovers kubernetes API preferred resources and returns them
func (k *Kraph) discoverAPI(ctx context.Context) (*API, error) {
	srvPrefResList, err := k.disc.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch api groups from kubernetes: %w", err)
	}

	api := &API{
		resourceMap: make(map[string][]Resource),
	}

	for _, srvPrefRes := range srvPrefResList {
		gv, err := schema.ParseGroupVersion(srvPrefRes.GroupVersion)
		if err != nil {
			return nil, fmt.Errorf("failed parsing %s into GroupVersion: %w", srvPrefRes.GroupVersion, err)
		}

		for _, apiResource := range srvPrefRes.APIResources {
			if !provides(apiResource.Verbs, "list") {
				continue
			}

			resource := Resource{
				r:  apiResource,
				gv: gv,
			}

			api.resources = append(api.resources, resource)
			for _, path := range resource.Paths() {
				api.resourceMap[path] = append(api.resourceMap[path], resource)
			}
		}
	}

	return api, nil
}
