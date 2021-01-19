package graph

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/entity"
	"github.com/milosgajdos/kraph/pkg/query"
	"gonum.org/v1/gonum/graph/encoding"
)

// Entity is graph entity
type Entity interface {
	entity.Entity
}

// DOTer implements GraphViz DOT properties.
type DOTer interface {
	// DOTID returns Graphviz DOT ID.
	DOTID() string
	// SetDOTID sets Graphviz DOT ID.
	SetDOTID(string)
}

// DOTNode is a GraphViz DOT Node.
type DOTNode interface {
	DOTer
	Node
}

// DOTEdge is a GraphViz DOT Edge.
type DOTEdge interface {
	DOTer
	Edge
}

// Node is a graph node.
type Node interface {
	Entity
	// Object returns api.Object.
	Object() api.Object
}

// Edge is an edge between two graph nodes.
type Edge interface {
	Entity
	// FromNode returns the from node of the edge.
	FromNode() Node
	// ToNode returns the to node of the edge.
	ToNode() Node
	// Weight returns edge weight.
	Weight() float64
}

// DOTGraph returns GraphViz DOT graph.
type DOTGraph interface {
	Graph
	// DOTID returns grapph DOT ID.
	DOTID() string
	// DOTAttributers returns graph DOT attributes.
	DOTAttributers() (graph, node, edge encoding.Attributer)
	// DOT returns Graphviz DOT graph.
	DOT() (string, error)
}

// Graph is a graph of API objects.
type Graph interface {
	// NewNode creates new node and returns it
	// TODO: clean up this interface
	NewNode(api.Object, ...entity.Option) (Node, error)
	// AddNode adds a new node to the graph.
	AddNode(Node) error
	// Node returns node with the given id.
	Node(id string) (Node, error)
	// Nodes returns all graph nodes.
	Nodes() ([]Node, error)
	// RemoveNode removes node from the graph.
	RemoveNode(id string) error
	// Link links two nodes and returns the new edge.
	Link(from, to string, opts LinkOptions) (Edge, error)
	// Edge returns the edge between the two nodes.
	Edge(from, to string) (Edge, error)
	// Edges returns all graph edges.
	Edges() ([]Edge, error)
	// RemoveEdge removes edge from the graph.
	RemoveEdge(from, to string) error
	// SubGraph returns a subgraph starting at node
	// with the given id up to the given depth.
	SubGraph(id string, depth int) (Graph, error)
	// Query the graph and return the results.
	Query(*query.Query) ([]Entity, error)
}
