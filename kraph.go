package kraph

import (
	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/store"
	"github.com/milosgajdos/kraph/store/memory"
)

// Kraph builds a graph of API objects
type Kraph interface {
	// Build builds a graph and returns graph store
	Build(api.Client) (store.Graph, error)
	// Store returns graph store
	Store() store.Store
}

// Options are kraph options
type Options struct {
	Store store.Store
}

// Option is functional kraph option
type Option func(*Options)

// Store configures kraph store
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// NewOptions creates default options and returns it
func NewOptions() (*Options, error) {
	m, err := memory.NewStore("default")
	if err != nil {
		return nil, err
	}

	return &Options{
		Store: m,
	}, nil
}
