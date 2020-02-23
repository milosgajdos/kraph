package kraph

import (
	"errors"
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

var (
	ErrNotImplemented = errors.New("not implemented")
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
	// client is kubernetes clientset
	client kubernetes.Interface
	// Global DOT attributes
	GraphAttrs Attributes
	NodeAttrs  Attributes
	EdgeAttrs  Attributes
}

// New creates new Kraph and returns it
func New(client kubernetes.Interface) (*Kraph, error) {
	return &Kraph{
		WeightedUndirectedGraph: simple.NewWeightedUndirectedGraph(0.0, 0.0),
		nodeMap:                 make(map[types.UID]*Node),
		client:                  client,
	}, nil
}

// NewNode creates new kraph node and returns it
func (k *Kraph) NewNode(name string) *Node {
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

// Build builds the kubernetes resource graph
func (k *Kraph) Build() error {
	if err := k.addNodes(); err != nil {
		return fmt.Errorf("failed building kraph: %w", err)
	}

	return nil
}

// addNodes gets a list of kubernetes nodes and adds them to kraph
func (k *Kraph) addNodes() error {
	// simple options for now
	options := metav1.ListOptions{
		Limit: 100,
	}

	nodes, err := k.client.CoreV1().Nodes().List(options)
	if err != nil {
		return fmt.Errorf("failed getting nodes: %v", err)
	}

	// iterate through nodes and add them to graph
	for _, node := range nodes.Items {
		fmt.Printf("Adding node %v to graph", node)
		if n, ok := k.nodeMap[node.UID]; ok {
			if knode := k.Node(n.ID()); knode == nil {
				k.AddNode(k.NewNode(node.Name))
			}
			continue
		}

		knode := k.NewNode(node.Name)
		k.nodeMap[node.UID] = knode
		k.AddNode(knode)
	}

	return nil
}

// addNamespaces gets a list of all kubernetes namespaces and adds them to kraph
func (k *Kraph) addNamespaces() error {
	return ErrNotImplemented
}
