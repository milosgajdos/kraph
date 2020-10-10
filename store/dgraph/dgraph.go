package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	//dgo "github.com/dgraph-io/dgo/v200"
	dgapi "github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/google/uuid"
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/attrs"
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
			uid
			xid
			name
			kind
			namespace
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

	attrs := attrs.New()
	attrs.Set("name", n.Name)
	attrs.Set("kind", n.Kind)
	attrs.Set("namespace", n.Namespace)

	node := entity.NewNode(n.XID, entity.Attrs(attrs))

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
			namespace
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
		attrs := attrs.New()
		attrs.Set("name", n.Name)
		attrs.Set("kind", n.Kind)
		attrs.Set("namespace", n.Namespace)

		node := entity.NewNode(n.XID, entity.Attrs(attrs))

		nodes[i] = node
	}

	return nodes, nil
}

// Edge returns the edge from uid to vid if it exists and nil otherwise
func (d *dgraph) Edge(uid, vid string) (store.Edge, error) {
	q := `
	query Edge($uid: string, $vid: string) {
	  node(func: eq(xid, $uid)) @cascade {
		uid
		dlink as link @filter(eq(xid, $vid)) @facets(drel as relation, dweight as weight, dluid as luid) {
			uid
		}
		rlink as ~link @filter(eq(xid, $vid)) @facets(rrel as relation, rweight as weight, rluid as luid) {
			uid
		}
	  }

	  fromUid(func: uid(dlink)) {
		xid
		relation: val(drel)
		weight: val(dweight)
		luid: val(dluid)
	  }

	   fromVid(func: uid(rlink)) {
		xid
		relation: val(rrel)
		weight: val(rweight)
		luid: val(rluid)
	  }
	}
	`

	var r struct {
		FromUid []Node `json:"fromUid"`
		FromVid []Node `json:"fromVid"`
	}

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$uid": uid, "$vid": vid})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNodeNotFound, err)
	}

	if err = json.Unmarshal(resp.Json, &r); err != nil {
		return nil, err
	}

	fmt.Println(resp)

	var weight float64
	var relation string
	var luid string

	switch {
	case len(r.FromUid) > 0:
		if rel := r.FromUid[0].Relation; rel != nil {
			relation = *rel
		}
		if w := r.FromUid[0].Weight; w != nil {
			weight = *w
		}
		if l := r.FromUid[0].LUID; l != nil {
			luid = *l
		}
	case len(r.FromVid) > 0:
		if rel := r.FromVid[0].Relation; rel != nil {
			relation = *rel
		}
		if w := r.FromVid[0].Weight; w != nil {
			weight = *w
		}
		if l := r.FromVid[0].LUID; l != nil {
			luid = *l
		}
	default:
		return nil, errors.ErrEdgeNotExist
	}

	if luid == "" {
		luid = uuid.New().String()
	}
	edge := entity.NewEdge(luid, entity.NewNode(uid), entity.NewNode(vid), entity.Relation(relation), entity.Weight(weight))

	return edge, nil
}

// Add adds an API object to the dgraph store and returns it
func (d *dgraph) Add(obj api.Object, opts store.AddOptions) (store.Entity, error) {
	query := `
	{
		node(func: eq(xid, "` + obj.UID().String() + `")) {
			u as uid
		}
	}
	`

	node := &Node{
		UID:       "uid(u)",
		XID:       obj.UID().String(),
		Name:      obj.Kind() + "-" + obj.Name(),
		Kind:      obj.Kind(),
		Namespace: obj.Namespace(),
		CreatedAt: time.Now(),
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
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	if _, err := txn.Do(ctx, req); err != nil {
		return nil, err
	}

	snode := entity.NewNode(obj.UID().String())

	return snode, nil
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edge between the nodes already exists.
// It returns error if either of the nodes does not exist in the graph.
func (d *dgraph) Link(from store.Node, to store.Node, opts store.LinkOptions) (store.Edge, error) {
	query := `
	{
		from as var(func: eq(xid, "` + from.UID() + `")) {
			fid as uid
		}

		to as var(func: eq(xid, "` + to.UID() + `")) {
			tid as uid
		}
	}
	`

	toNode := &Node{
		UID:   "uid(tid)",
		DType: []string{"Object"},
		LUID:  DgString(uuid.New().String()),
	}

	if opts.Relation != "" {
		toNode.Relation = DgString(opts.Relation)
	}

	if opts.Weight > 0.0 {
		toNode.Weight = DgFloat(opts.Weight)
	}

	node := &Node{
		UID:   "uid(fid)",
		DType: []string{"Object"},
		Link:  []Node{*toNode},
	}

	pb, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	mu := &dgapi.Mutation{
		Cond:    `@if(gt(len(from), 0) AND gt(len(to), 0))`,
		SetJson: pb,
	}

	req := &dgapi.Request{
		Query:     query,
		Mutations: []*dgapi.Mutation{mu},
		CommitNow: true,
	}

	ctx := context.Background()
	txn := d.client.NewTxn()
	defer txn.Discard(ctx)

	if _, err := txn.Do(ctx, req); err != nil {
		return nil, err
	}

	edge := entity.NewEdge(*toNode.LUID, from, to, entity.Relation(opts.Relation), entity.Weight(opts.Weight))

	return edge, nil
}

// Delete deletes an entity from the store.
func (d *dgraph) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case store.Node:
		q := `
		query Node($xid: string) {
			node(func: eq(xid, $xid)) {
				uid
				xid
			}
		 }
		`

		ctx := context.Background()
		txn := d.client.NewTxn()
		defer txn.Discard(ctx)

		resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$xid": v.UID()})
		if err != nil {
			return err
		}

		var r struct {
			Result []Node `json:"node"`
		}

		if err = json.Unmarshal(resp.Json, &r); err != nil {
			return err
		}

		res := len(r.Result)

		switch {
		case res == 0:
			return errors.ErrNodeNotFound
		case res > 1:
			return errors.ErrDuplicateNode
		}

		n := r.Result[0]

		node := map[string]string{"uid": n.UID}
		pb, err := json.Marshal(node)
		if err != nil {
			return err
		}

		mu := &dgapi.Mutation{
			CommitNow:  true,
			DeleteJson: pb,
		}

		ctx = context.Background()
		_, err = d.client.NewTxn().Mutate(ctx, mu)
		if err != nil {
			return err
		}
	case store.Edge:
		query := `
		{
			from as var(func: eq(xid, "` + v.From().UID() + `")) {
				fid as uid
			}

			to as var(func: eq(xid, "` + v.To().UID() + `")) {
				tid as uid
			}
		}
		`

		node := &Node{
			UID:   "uid(fid)",
			DType: []string{"Object"},
			Link: []Node{
				{UID: "uid(tid)", DType: []string{"Object"}},
			},
		}

		pb, err := json.Marshal(node)
		if err != nil {
			return err
		}

		mu := &dgapi.Mutation{
			Cond:       `@if(gt(len(from), 0) AND gt(len(to), 0))`,
			DeleteJson: pb,
		}

		req := &dgapi.Request{
			Query:     query,
			Mutations: []*dgapi.Mutation{mu},
			CommitNow: true,
		}

		ctx := context.Background()
		txn := d.client.NewTxn()
		defer txn.Discard(ctx)

		if _, err := txn.Do(ctx, req); err != nil {
			return err
		}
	default:
		return errors.ErrUnknownEntity
	}

	return nil
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
func (d *dgraph) SubGraph(n store.Node, depth int) (store.Graph, error) {
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
