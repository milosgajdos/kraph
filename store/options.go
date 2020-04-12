package store

var (
	// DefaultEdgeWeight defines default edge weight
	DefaultEdgeWeight = 0.0
)

// Options are store options
type Options struct {
	Metadata   Metadata
	Attributes Attributes
	GraphAttrs Attributes
	NodeAttrs  Attributes
	EdgeAttrs  Attributes
	Weight     float64
}

// Option sets options
type Option func(*Options)

// Meta sets entity metadata
func Meta(m Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// Attrs sets entity attributes
func Attrs(a Attributes) Option {
	return func(o *Options) {
		o.Attributes = a
	}
}

// GraphAttrs sets graph attributes
func GraphAttrs(a Attributes) Option {
	return func(o *Options) {
		o.GraphAttrs = a
	}
}

// NodeAttrs sets global node attributes
func NodeAttrs(a Attributes) Option {
	return func(o *Options) {
		o.NodeAttrs = a
	}
}

// EdgeAttrs sets global edge attributes
func EdgeAttrs(a Attributes) Option {
	return func(o *Options) {
		o.EdgeAttrs = a
	}
}

// Weight sets weight option
func Weight(w float64) Option {
	return func(o *Options) {
		o.Weight = w
	}
}

// NewOptions returns empty options
func NewOptions() Options {
	return Options{
		Metadata:   NewMetadata(),
		Attributes: NewAttributes(),
		GraphAttrs: NewAttributes(),
		NodeAttrs:  NewAttributes(),
		EdgeAttrs:  NewAttributes(),
		Weight:     DefaultEdgeWeight,
	}
}
