package store

import (
	"github.com/milosgajdos/kraph/store/attrs"
	"github.com/milosgajdos/kraph/store/metadata"
)

const (
	// DefaultWeight is default weight
	DefaultWeight = 1.0
)

// Options are store options
type Options struct{}

// Option configures store
type Option func(*Options)

func NewOptions() Options {
	return Options{}
}

// AddOptions are store options
type AddOptions struct {
	Attrs    *attrs.Attrs
	Metadata *metadata.Metadata
}

// AddOption sets options
type AddOption func(*AddOptions)

// NewOptions returns empty options
func NewAddOptions() AddOptions {
	return AddOptions{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
	}
}

// DelOptions are store options
type DelOptions struct {
	Attrs    *attrs.Attrs
	Metadata *metadata.Metadata
}

// DelOption sets options
type DelOption func(*DelOptions)

// NewDelOptions returns empty options
func NewDelOptions() DelOptions {
	return DelOptions{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
	}
}

// LinkOptions are store options
type LinkOptions struct {
	Attrs    *attrs.Attrs
	Metadata *metadata.Metadata
	Weight   float64
	Relation string
}

// LinkOption sets options
type LinkOption func(*LinkOptions)

// NewLinkOptions returns empty options
func NewLinkOptions() LinkOptions {
	return LinkOptions{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
		Weight:   DefaultWeight,
	}
}
