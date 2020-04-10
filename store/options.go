package store

// Options are store options
type Options struct {
	Metadata   Metadata
	Attributes Attributes
	Weight     float64
}

// Option sets options
type Option func(*Options)

// Attrs sets entity attributes
func Attrs(a Attributes) Option {
	return func(o *Options) {
		o.Attributes = a
	}
}

// Meta sets entity metadata
func Meta(m Metadata) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

// Weight sets weight option
func Weight(w float64) Option {
	return func(o *Options) {
		o.Weight = w
	}
}
