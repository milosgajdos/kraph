package graph

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
)

const (
	// DefaultWeight is default edge weight
	DefaultWeight = 1.0
)

// DOTOptions are DOT graph options.
type DOTOptions struct {
	GraphAttrs attrs.DOT
	NodeAttrs  attrs.DOT
	EdgeAttrs  attrs.DOT
}

// DOTOption configures DOT graph.
type DOTOption func(*Options)

// Options are graph options.
type Options struct {
	DOTOptions DOTOptions
	Weight     float64
}

// NewOptions returns default graph options.
func NewOptions() Options {
	return Options{
		Weight: DefaultWeight,
	}
}

// LinkOptions are graph link options.
type LinkOptions struct {
	Attrs  attrs.Attrs
	Weight float64
}

// LinkOption sets link options.
type LinkOption func(*LinkOptions)

// NewLinkOptions returns default link options.
func NewLinkOptions() LinkOptions {
	return LinkOptions{
		Weight: DefaultWeight,
		Attrs:  attrs.New(),
	}
}
