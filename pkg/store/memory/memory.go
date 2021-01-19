package memory

import (
	"github.com/milosgajdos/kraph/pkg/graph"
	"github.com/milosgajdos/kraph/pkg/graph/memory"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/store"
)

// Memory is in-memory store.
type Memory struct {
	// id is the store id
	id string
	// g is the store graph
	g graph.Graph
	// options are the store options
	opts store.Options
}

// NewStore creates a new in-memory store and returns it.
func NewStore(id string, opts store.Options) (*Memory, error) {
	g := opts.Graph

	if g == nil {
		var err error
		g, err = memory.NewWUG(id+"-graph", graph.Options{})
		if err != nil {
			return nil, err
		}
	}

	return &Memory{
		id:   id,
		g:    g,
		opts: opts,
	}, nil
}

// ID returns store ID.
func (m Memory) ID() string {
	return m.id
}

// Options returns memory store options.
func (m Memory) Options() store.Options {
	return m.opts
}

// Graph returns graph handle.
func (m *Memory) Graph() graph.Graph {
	return m.g
}

// Add stores entity in memory store.
func (m *Memory) Add(e store.Entity, opts store.AddOptions) error {
	switch v := e.(type) {
	case graph.Node:
		return m.g.AddNode(v)
	case graph.Edge:
		from := v.FromNode().UID()
		to := v.ToNode().UID()

		if _, err := m.g.Link(from, to, graph.LinkOptions{Attrs: opts.Attrs}); err != nil {
			return err
		}
		return nil
	}

	return store.ErrUnknownEntity
}

// Delete deletes entity e from memory store.
func (m *Memory) Delete(e store.Entity, opts store.DelOptions) error {
	switch v := e.(type) {
	case graph.Node:
		return m.g.RemoveNode(v.UID())
	case graph.Edge:
		return m.g.RemoveEdge(v.FromNode().UID(), v.ToNode().UID())
	}

	return store.ErrUnknownEntity
}

// Query queries the store and returns the results
func (m Memory) Query(q *query.Query) ([]store.Entity, error) {
	qents, err := m.g.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]store.Entity, len(qents))

	for i, e := range qents {
		results[i] = e.(store.Entity)
	}

	return results, nil
}
