package kraph

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/store"
)

// Filter lets you skip adding api.Object to kraph.
type Filter func(api.Object) bool

// Kraph builds a graph of API objects.
type Kraph interface {
	// Build builds a graph of an API
	Build(api.Scraper, ...Filter) error
	// Store returns graph store.
	Store() store.Store
	// Netadata returns kraph metadata.
	Metadata() metadata.Metadata
}

// Options are kraph options.
type Options struct {
	Metadata metadata.Metadata
}

// Option is functional kraph option.
type Option func(*Options)

// Metadata configures krap metadata.
func Metadata(m metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// NewOptions creates default options and returns it.
func NewOptions() (*Options, error) {
	return &Options{
		Metadata: metadata.New(),
	}, nil
}
