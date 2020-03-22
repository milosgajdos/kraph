package query

// Options are query options
type Options struct {
	Namespace string
	Kind      string
	Name      string
	Attrs     map[string]string
}

// Option configures the query
type Option func(*Options)

// Namespace configures options namespace
func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

// Kind configures options kind
func Kind(kind string) Option {
	return func(o *Options) {
		o.Kind = kind
	}
}

// Name configures options name
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Attrs configures attributes
func Attrs(a map[string]string) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// NewOptions returns default options
func NewOptions() Options {
	return Options{
		Namespace: "",
		Kind:      "",
		Name:      "",
		Attrs:     make(map[string]string),
	}
}
