package entity

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
)

// Options are store options
type Options struct {
	Attrs attrs.Attrs
}

// Option sets options
type Option func(*Options)

// NewOptions returns empty options
func NewOptions() Options {
	return Options{
		Attrs: attrs.New(),
	}
}

// Attrs sets entity attributes
func Attrs(a attrs.Attrs) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}
