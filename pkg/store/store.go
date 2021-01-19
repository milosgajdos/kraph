package store

import (
	"github.com/milosgajdos/kraph/pkg/graph"
	"github.com/milosgajdos/kraph/pkg/query"
)

// Entity is store entity.
type Entity interface {
	graph.Entity
}

// Store allows to store and query entities.
type Store interface {
	// Graph returns the graph handle.
	Graph() graph.Graph
	// Add entity to the store.
	Add(Entity, AddOptions) error
	// Delete entity from the store.
	Delete(Entity, DelOptions) error
	// Query the store and return the results.
	Query(*query.Query) ([]Entity, error)
}
