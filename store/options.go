package store

var (
	// DefaultEdgeWeight defines default edge weight
	DefaultEdgeWeight = 0.0
)

// Options are store options
type Options struct {
	Metadata   Metadata
	EntAttrs   Attrs
	GraphAttrs Attrs
	NodeAttrs  Attrs
	EdgeAttrs  Attrs
	Weight     float64
	Relation   string
}

// Option sets options
type Option func(*Options)

// Meta sets entity metadata
func Meta(m Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// EntAttrs sets entity attributes
func EntAttrs(a Attrs) Option {
	return func(o *Options) {
		o.EntAttrs = a
	}
}

// GraphAttrs sets graph attributes
func GraphAttrs(a Attrs) Option {
	return func(o *Options) {
		o.GraphAttrs = a
	}
}

// NodeAttrs sets global node attributes
func NodeAttrs(a Attrs) Option {
	return func(o *Options) {
		o.NodeAttrs = a
	}
}

// EdgeAttrs sets global edge attributes
func EdgeAttrs(a Attrs) Option {
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

// Relation sets relation option
func Relation(r string) Option {
	return func(o *Options) {
		o.Relation = r
	}
}

// NewOptions returns empty options
func NewOptions() Options {
	return Options{
		Metadata:   NewMetadata(),
		EntAttrs:   NewAttributes(),
		GraphAttrs: NewAttributes(),
		NodeAttrs:  NewAttributes(),
		EdgeAttrs:  NewAttributes(),
		Weight:     DefaultEdgeWeight,
		Relation:   "link",
	}
}
