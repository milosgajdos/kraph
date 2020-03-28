package query

// Options are query options
type Options struct {
	Namespace string
	Kind      string
	Name      string
	Version   string
	Group     string
	Attrs     map[string]string
}

// Option configures the query
type Option func(*Options)

// Namespace configures namespace option
func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

// Kind configures kind option
func Kind(kind string) Option {
	return func(o *Options) {
		o.Kind = kind
	}
}

// Name configures name option
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Attrs configures attributes option
func Attrs(a map[string]string) Option {
	return func(o *Options) {
		o.Attrs = a
	}
}

// Version configues version option
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Group configures group option
func Group(g string) Option {
	return func(o *Options) {
		o.Group = g
	}
}

// NewOptions returns default options
func NewOptions() Options {
	return Options{
		Namespace: "",
		Kind:      "",
		Name:      "",
		Version:   "",
		Group:     "",
		Attrs:     make(map[string]string),
	}
}
