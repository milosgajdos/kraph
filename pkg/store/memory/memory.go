package memory

import (
	"fmt"
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/errors"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/entity"
	"github.com/milosgajdos/kraph/pkg/uuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

// Memory is in-memory graph store
type Memory struct {
	// g is store graph
	g *Graph
	// id is the store id
	id string
	// options are store options
	opts store.Options
}

// NewStore creates new in-memory store and returns it
func NewStore(id string, opts store.Options) (*Memory, error) {
	return &Memory{
		g:    NewGraph(id, opts.GraphOptions),
		id:   id,
		opts: opts,
	}, nil
}

// ID returns store ID
func (m Memory) ID() string {
	return m.id
}

// Options returns store options
func (m Memory) Options() store.Options {
	return m.opts
}

// Graph returns graph handle
func (m *Memory) Graph() store.Graph {
	return m.g
}

// Add adds obj to the store and returns it
func (m *Memory) Add(obj api.Object, opts store.AddOptions) (store.Entity, error) {
	return m.g.NewNode(obj, opts)
}

// Delete deletes entity e from the memory store
func (m *Memory) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case store.Edge:
		return m.g.RemoveLine(v.From().UID(), v.To().UID(), v.UID())
	case store.Node:
		return m.g.RemoveNode(v.UID())
	default:
		return errors.ErrUnknownEntity
	}
}

// QueryNode returns all the nodes that match given query.
func (m *Memory) QueryNode(q *query.Query) ([]*Node, error) {
	match := q.Matcher()

	if quid := match.UID(); quid != nil {
		if uid, ok := quid.Value().(uuid.UID); ok && len(uid.String()) > 0 {
			if n, ok := m.g.nodes[uid.String()]; ok {
				return []*Node{n}, nil
			}
		}
	}

	var results []*Node

	visit := func(n graph.Node) {
		node := n.(*Node)
		nodeObj := node.Metadata().Get("object").(api.Object)

		if match.NamespaceVal(nodeObj.Namespace()) {
			if match.KindVal(nodeObj.Resource().Kind()) {
				if match.NameVal(nodeObj.Name()) {
					if !match.AttrsVal(node.Attrs()) {
						return
					}

					// create a deep copy of the matched node
					attrs := attrs.New()
					metadata := metadata.New()

					for _, k := range node.Attrs().Keys() {
						attrs.Set(k, node.Attrs().Get(k))
					}

					for _, k := range node.Metadata().Keys() {
						metadata.Set(k, node.Metadata().Get(k))
					}

					dotid := strings.Join([]string{
						nodeObj.Resource().Version(),
						nodeObj.Namespace(),
						nodeObj.Resource().Kind(),
						nodeObj.Name()}, "/")
					attrs.Set("name", dotid)

					entOpts := []entity.Option{
						entity.Metadata(metadata),
						entity.Attrs(attrs),
					}

					n := NewNode(node.ID(), node.UID(), dotid, entOpts...)

					results = append(results, n)
				}
			}
		}
	}

	dfs := traverse.DepthFirst{
		Visit: visit,
	}

	dfs.WalkAll(m.g.WeightedUndirectedGraph, nil, nil, func(graph.Node) {})

	return results, nil
}

// QueryEdge returns all the edges that match given query
func (m *Memory) QueryLine(q *query.Query) ([]*Line, error) {
	match := q.Matcher()

	var results []*Line

	if quid := match.UID(); quid != nil {
		if uid, ok := quid.Value().(string); ok && len(uid) > 0 {
			if l, ok := m.g.lines[uid]; ok {
				return []*Line{l}, nil
			}
		}
	}

	trav := func(e graph.Edge) bool {
		from := e.From().(*Node)
		to := e.To().(*Node)

		if lines := m.g.WeightedLines(from.ID(), to.ID()); lines != nil {
			for lines.Next() {
				wl := lines.WeightedLine()
				we := wl.(*Line).Edge
				if match.WeightVal(we.Weight()) {
					if !match.AttrsVal(we.Attrs()) {
						continue
					}

					attrs := attrs.New()
					metadata := metadata.New()

					for _, k := range we.Attrs().Keys() {
						attrs.Set(k, we.Attrs().Get(k))
					}

					for _, k := range we.Metadata().Keys() {
						metadata.Set(k, we.Metadata().Get(k))
					}

					opts := []entity.Option{
						entity.Attrs(attrs),
						entity.Metadata(metadata),
						entity.Weight(we.Weight()),
					}

					ent := NewLine(wl.ID(), we.UID(), we.UID(), from, to, opts...)

					results = append(results, ent)
				}
			}
		}

		return true
	}

	dfs := traverse.DepthFirst{
		Traverse: trav,
	}

	dfs.WalkAll(m.g.WeightedUndirectedGraph, nil, nil, func(graph.Node) {})

	return results, nil
}

// Query queries the in-memory graph and returns the matched results.
func (m *Memory) Query(q *query.Query) ([]store.Entity, error) {
	var e query.Entity

	if m := q.Matcher().Entity(); m != nil {
		var ok bool
		e, ok = m.Value().(query.Entity)
		if !ok {
			return nil, errors.ErrInvalidEntity
		}
	}

	var entities []store.Entity

	switch e {
	case query.Node:
		nodes, err := m.QueryNode(q)
		if err != nil {
			return nil, fmt.Errorf("Node query: %w", err)
		}
		for _, node := range nodes {
			entities = append(entities, node.Node)
		}
	case query.Edge:
		edges, err := m.QueryLine(q)
		if err != nil {
			return nil, fmt.Errorf("Edge query: %w", err)
		}
		for _, edge := range edges {
			entities = append(entities, edge)
		}
	default:
		return nil, errors.ErrUnknownEntity
	}

	return entities, nil
}
