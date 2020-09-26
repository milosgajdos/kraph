package store

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	"gonum.org/v1/gonum/graph/encoding"
)

// Metadata is a simple key-value store
// for storing arbitrary metadata.
type Metadata interface {
	// Keys returns all metadata keys
	Keys() []string
	// Get returns the metadata for the given key
	Get(string) interface{}
	// Set stores the metadata
	Set(string, interface{})
}

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

// Entity is an arbitrary store entity
type Entity interface {
	// ID returns node ID
	ID() string
	// Name is the name of the entity
	Name() string
	// Attrs returns entity attributes
	Attrs() Attrs
	// Metadata returns entity metadata
	Metadata() Metadata
	// Attributes returns attributes as a slice of encoding.Attribute
	Attributes() []encoding.Attribute
}

// DOTNode is a GraphViz DOT Node
type DOTNode interface {
	Node
	// DOTID returns Graphiz DOT ID
	DOTID() string
	// SetDOTID sets Graphiz DOT ID
	SetDOTID(string)
}

// Node is a graph node
type Node interface {
	Entity
}

// Edge is an edge between two graph nodes
type Edge interface {
	Entity
	// From returns the from node of the edge
	From() Node
	// To returns the to node of the edge.
	To() Node
	// Weight returns edge weight
	Weight() float64
}

// DOTGraph returns Graphiz DOT graph
type DOTGraph interface {
	Graph
	// DOTID returns graph DOT ID
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
	// Edge returns the edge from u to v, with IDs uid and vid,
	// if such an edge exists and nil otherwise
	Edge(uid, vid string) (Edge, error)
	// Subgraph returns a subgraph of the graph starting at Node
	// up to the given depth or it returns an error
	SubGraph(id string, depth int) (Graph, error)
}

// Store allows to store and query the graph of API objects
type Store interface {
	// Graph returns graph handle
	Graph() Graph
	// Add adds an api.Object to the store and returns a Node
	Add(api.Object, ...Option) (Node, error)
	// Link links two nodes and returns the new edge between them
	Link(Node, Node, ...Option) (Edge, error)
	// Delete deletes an entity from the store
	Delete(Entity, ...Option) error
	// Query queries the store and returns the results
	Query(...query.Option) ([]Entity, error)
}
