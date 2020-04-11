package memory

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/errors"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
)

// Memory is in-memory graph store
type Memory struct{}

// Add adds an API object to the in-memory graph as a graph node and returns it
func (m *Memory) Add(obj api.Object, opts ...store.Option) (store.Node, error) {
	return nil, errors.ErrNotImplemented
}

// Link creates a new edge between the nodes and returns it or it returns
// an existing edge if the edges between the nodes already exists.
func (m *Memory) Link(from store.Node, to store.Node, opts ...store.Option) (store.Edge, error) {
	return nil, errors.ErrNotImplemented
}

// Query queries the in-memory graph and returns the matched results.
func (m *Memory) Query(q ...query.Option) ([]store.Entity, error) {
	return nil, errors.ErrNotImplemented
}
