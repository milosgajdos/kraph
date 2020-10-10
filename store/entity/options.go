package entity

import (
	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/metadata"
)

const (
	// DefaultWeight is default weight
	DefaultWeight = 1.0
)

// Options are store options
type Options struct {
	Attrs    *attrs.Attrs
	Metadata *metadata.Metadata
	Weight   float64
	Relation string
}

// Option sets options
type Option func(*Options)

// NewOptions returns empty options
func NewOptions() Options {
	return Options{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
		Weight:   DefaultWeight,
	}
}

// Metadata sets entity metadata
func Metadata(m *metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// Attrs sets entity attributes
func Attrs(a *attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// Weight returns entity weight
func Weight(w float64) Option {
	return func(o *Options) {
		o.Weight = w
	}
}

// Relation configures entity relation
func Relation(r string) Option {
	return func(o *Options) {
		o.Relation = r
	}
}
