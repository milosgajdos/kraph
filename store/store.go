package store

import (
	"github.com/milosgajdos/kraph/query"
	"gonum.org/v1/gonum/graph/encoding"
)

// Attrs provide a simple key-value store
// for storing arbitrary entity attributes
type Attrs interface {
	// Keys returns all attribute keys
	Keys() []string
	// Get returns the attribute value for the given key
	Get(string) string
	// Set sets the value of the attribute for the given key
	Set(string, string)
}

// DOTAttrs are Attrs which implement graph.DOTAttributes
type DOTAttrs interface {
	// Attributes returns attributes as a slice of encoding.Attribute
	Attributes() []encoding.Attribute
}

// Metadata provides a simple key-valule store
// for arbitrary entity data of arbitrary type
type Metadata interface {
	// Keys returns all metadata keys
	Keys() []string
	// Get returns the attribute value for the given key
	Get(string) interface{}
	// Set sets the value of the attribute for the given key
	Set(string, interface{})
}

// Entity is an arbitrary store entity
type Entity interface {
	// ID returns unique ID
	ID() string
	// Attrs returns attributes
	Attrs() Attrs
	// Metadata returns metadata
	Metadata() Metadata
}

// DOTNode is a GraphViz DOT Node
type DOTNode interface {
	Node
	// DOTID returns Graphviz DOT ID
	DOTID() string
	// SetDOTID sets Graphviz DOT ID
	SetDOTID(string)
}

// Node is a graph node
type Node interface {
	Entity
}

// WeightedEdge is an edge with weight
type WeightedEdge interface {
	Edge
	// Weight returns edge weight
	Weight() float64
}

// Edge is an edge between two nodes
type Edge interface {
	Entity
	// From returns the from node of the edge
	From() Node
	// To returns the to node of the edge.
	To() Node
}

// DOTGraph returns Graphiz DOT store
type DOTGraph interface {
	Graph
	// DOTID returns DOT graph ID
	DOTID() string
	// DOTAttributers returns global graph DOT attributes
	DOTAttributers() (graph, node, edge encoding.Attributer)
	// DOT returns Graphviz graph
	DOT() (string, error)
}

// Graph is a graph of API objects
type Graph interface {
	// Node returns the node with the given ID if it exists
	// in the graph, and nil otherwise.
	Node(id string) (Node, error)
	// Nodes returns all the nodes in the graph.
	Nodes() ([]Node, error)
	// Link links two nodes and returns the new edge between them
	Link(Node, Node, ...Option) (Edge, error)
	// Edge returns the edge from u to v, with IDs uid and vid,
	// if such an edge exists and nil otherwise
	Edge(uid, vid string) (Edge, error)
	// Subgraph returns a subgraph of the graph starting at Node
	// up to the given depth or it returns an error
	SubGraph(Node, depth int) (Graph, error)
}

// Store allows to store and query the graph of API objects
type Store interface {
	// Add adds an api.Object to the store and returns a Node
	Add(Entity, ...Option) (Entity, error)
	// Delete deletes an entity from the store
	Delete(Entity, ...Option) error
	// Query queries the store and returns the results
	Query(...query.Option) ([]Entity, error)
}
