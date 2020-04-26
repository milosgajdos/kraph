package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	//dgo "github.com/dgraph-io/dgo/v200"
	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/entity"
	"gonum.org/v1/gonum/graph/encoding"
)

var (
	// DefaultURL is default dgraph store URL
	DefaultURL = "localhost:9080"
)

// dgraph is Dgraph store handle
type dgraph struct {
	id     string
	client *Client
}

// NewStore returns new dgraph store handle or error.
func NewStore(id string, client *Client, opts ...store.Option) (store.Store, error) {
	// NOTE: we do not use any options, yet
	// but we still initialize them nevertheless
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &dgraph{
		id:     id,
		client: client,
	}, nil
}

// Node returns the node with the given ID or error.
func (d *dgraph) Node(id string) (store.Node, error) {
	q := `
          query Node($xid: string){
		node(func: eq(xid, $xid)) {
			xid
			name
			kind
			ns
		}
          }
	`

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$xid": id})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNodeNotFound, err)
	}

	var r struct {
		Result []Node `json:"node"`
	}

	if err = json.Unmarshal(resp.Json, &r); err != nil {
		return nil, err
	}

	res := len(r.Result)

	switch {
	case res == 0:
		return nil, errors.ErrNodeNotFound
	case res > 1:
		return nil, errors.ErrDuplicateNode
	}

	n := r.Result[0]

	attrs := store.NewAttributes()
	attrs.Set("name", n.Name)
	attrs.Set("kind", n.Kind)
	attrs.Set("namespace", n.Namespace)

	node := entity.NewNode(n.UID, store.EntAttrs(attrs))

	return node, nil
}

// Nodes returns all the nodes in the graph.
func (d *dgraph) Nodes() ([]store.Node, error) {
	q := `
          query Nodes() {
		nodes(func: has(xid)) {
			xid
			name
			kind
			ns
		}
	  }
	`

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var r struct {
		Result []Node `json:"nodes"`
	}

	if err = json.Unmarshal(resp.Json, &r); err != nil {
		return nil, err
	}

	nodes := make([]store.Node, len(r.Result))

	for i, n := range r.Result {
		attrs := store.NewAttributes()
		attrs.Set("name", n.Name)
		attrs.Set("kind", n.Kind)
		attrs.Set("namespace", n.Namespace)

		node := entity.NewNode(n.UID, store.EntAttrs(attrs))

		nodes[i] = node
	}

	return nodes, nil
}

// Edge returns the edge from u to v, with IDs uid and vid
func (d *dgraph) Edge(uid, vid string) (store.Edge, error) {
	return nil, errors.ErrNotImplemented
}

// Add adds an API object to the dgraph store and returns it
func (d *dgraph) Add(obj api.Object, opts ...store.Option) (store.Node, error) {
	nodeOpts := store.NewOptions()
	for _, apply := range opts {
		apply(&nodeOpts)
	}

	query := `
          query Node($xid: string){
		node(func: eq(xid, $xid)) {
			xid
		}
          }
	`

	node := &Node{
		UID:       obj.UID().String(),
		Name:      obj.Kind() + "-" + obj.Name(),
		Kind:      obj.Kind(),
		Namespace: obj.Namespace(),
		DType:     []string{"Object"},
	}

	pb, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	mu := &dgapi.Mutation{
		SetJson: pb,
	}

	req := &dgapi.Request{
		Query:     query,
		Vars:      map[string]string{"$xid": node.UID},
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	if _, err := txn.Do(ctx, req); err != nil {
		return nil, err
	}

	snode := entity.NewNode(node.UID)

	return snode, nil
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
// It returns error if the edge failed to be added
// TODO: https://discuss.dgraph.io/t/dgraph-go-client-upsert-returning-uid/6148
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
func (d *dgraph) SubGraph(id string, depth int) (store.Graph, error) {
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
