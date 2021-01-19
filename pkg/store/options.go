package store

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/graph"
)

// Options are store options.
type Options struct {
	Graph graph.Graph
}

// AddOptions are store add options.
type AddOptions struct {
	Attrs attrs.Attrs
}

// AddOption sets add options.
type AddOption func(*AddOptions)

// NewOptions returns default add options.
func NewAddOptions() AddOptions {
	return AddOptions{
		Attrs: attrs.New(),
	}
}

// DelOptions are store delete options.
type DelOptions struct {
	Attrs attrs.Attrs
}

// DelOption sets delete options.
type DelOption func(*DelOptions)

// NewDelOptions returns default delete options.
func NewDelOptions() DelOptions {
	return DelOptions{
		Attrs: attrs.New(),
	}
}
