package dgraph

import (
	dgo "github.com/dgraph-io/dgo/v200"
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// dgraph is Dgraph DB handle
type dgraph struct {
	id     string
	client *dgo.Dgraph
}

// NewStore returns new dgraph store
func NewStore(id string, client *dgo.Dgraph, opts ...store.Option) (store.Store, error) {
	return &dgraph{
		id:     id,
		client: client,
	}, nil
}

// Node returns the node with the given ID if it exists
func (d *dgraph) Node(id int64) graph.Node {
	return nil
}

// Nodes returns all the nodes in the graph.
func (d *dgraph) Nodes() graph.Nodes {
	return nil
}

// From returns all nodes that can be reached directly
// from the node with the given ID.
func (d *dgraph) From(id int64) graph.Nodes {
	return nil
}

// HasEdgeBetween returns whether an edge exists between
// nodes with IDs xid and yid without considering direction.
func (d *dgraph) HasEdgeBetween(xid, yid int64) bool {
	return false
}

// Edge returns the edge from u to v, with IDs uid and vid
func (d *dgraph) Edge(uid, vid int64) graph.Edge {
	return nil
}

// Add adds an API object to the dgraph store and returns it
func (d *dgraph) Add(obj api.Object, opts ...store.Option) (store.Node, error) {
	return nil, errors.ErrNotImplemented
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It returns error if the edge failed to be added
func (d *dgraph) Link(from store.Node, to store.Node, opts ...store.Option) (store.Edge, error) {
	return nil, errors.ErrNotImplemented
}

// Delete deletes an entity from the memory store
func (d *dgraph) Delete(e store.Entity, opts ...store.Option) error {
	return errors.ErrNotImplemented
}

// QueryNode returns all the nodes that match given query.
func (d *dgraph) QueryNode(opts ...query.Option) ([]store.Node, error) {
	return nil, errors.ErrNotImplemented
}

// QueryEdge returns all the edges that match given query
func (d *dgraph) QueryEdge(opts ...query.Option) ([]store.Edge, error) {
	return nil, errors.ErrNotImplemented
}

// Query queries dgraph and returns the matched results.
func (d *dgraph) Query(q ...query.Option) ([]store.Entity, error) {
	return nil, errors.ErrNotImplemented
}

// SubGraph returns the subgraph of the node up to given depth or returns error
func (d *dgraph) SubGraph(n store.Node, depth int) (graph.Graph, error) {
	return nil, errors.ErrNotImplemented
}

// DOTID returns the store DOT ID.
func (d *dgraph) DOTID() string {
	return d.id
}

// DOTAttributers returns the global DOT graph attributers.
func (d *dgraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return nil, nil, nil
}

// DOT returns the GrapViz dot representation of kraph.
func (d *dgraph) DOT() (string, error) {
	return "", errors.ErrNotImplemented
}
