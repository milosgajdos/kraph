package store

var (
	// DefaultEdgeWeight defines default edge weight
	DefaultEdgeWeight = 0.0
)

// Options are store options
type Options struct {
	Metadata   Metadata
	Attributes Attrs
	Weight     float64
	Relation   string
}

// Option sets options
type Option func(*Options)

// NewOptions returns empty options
func NewOptions() Options {
	return Options{
		Metadata:   NewMetadata(),
		Attributes: NewAttributes(),
		Weight:     DefaultEdgeWeight,
		Relation:   "link",
	}
}

// Meta sets entity metadata
func Meta(m Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// Attributes sets entity attributes
func Attributes(a Attrs) Option {
	return func(o *Options) {
		o.Attributes = a
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
