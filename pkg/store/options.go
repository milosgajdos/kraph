package store

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/metadata"
)

const (
	// DefaultWeight is default weight
	DefaultWeight = 1.0
)

// DOTOptions are DOT options
type DOTOptions struct {
	GraphAttrs attrs.DOT
	NodeAttrs  attrs.DOT
	EdgeAttrs  attrs.DOT
}

// GraphOptions are graph options
type GraphOptions struct {
	DOTOptions DOTOptions
}

// DOTOption configures store
type DOTOption func(*Options)

// Options are store options
type Options struct {
	GraphOptions GraphOptions
}

// Option configures store
type Option func(*Options)

func NewOptions() Options {
	return Options{}
}

// AddOptions are store options
type AddOptions struct {
	Attrs    attrs.Attrs
	Metadata metadata.Metadata
}

// AddOption sets options
type AddOption func(*AddOptions)

// NewOptions returns default add options
func NewAddOptions() AddOptions {
	return AddOptions{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
	}
}

// DelOptions are store options
type DelOptions struct {
	Attrs    attrs.Attrs
	Metadata metadata.Metadata
}

// DelOption sets options
type DelOption func(*DelOptions)

// NewDelOptions returns default del options
func NewDelOptions() DelOptions {
	return DelOptions{
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
	}
}

// LinkOptions are link options
type LinkOptions struct {
	Line     bool
	Weight   float64
	Relation string
	Attrs    attrs.Attrs
	Metadata metadata.Metadata
}

// LinkOption sets options
type LinkOption func(*LinkOptions)

// NewLinkOptions returns default link options
func NewLinkOptions() LinkOptions {
	return LinkOptions{
		Weight:   DefaultWeight,
		Attrs:    attrs.New(),
		Metadata: metadata.New(),
	}
}
