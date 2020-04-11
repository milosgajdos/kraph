package query

// Options are query options
type Options struct {
	Namespace string
	Kind      string
	Name      string
	Version   string
	UID       string
	Group     string
	Weight    float64
	Entity    string
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

// UID configues uid option
func UID(u string) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// Group configures group option
func Group(g string) Option {
	return func(o *Options) {
		o.Group = g
	}
}

// Weight configures weight option
func Weight(w float64) Option {
	return func(o *Options) {
		o.Weight = w
	}
}

// Entity configures entity option
func Entity(e string) Option {
	return func(o *Options) {
		o.Entity = e
	}
}

// NewOptions returns default options
func NewOptions() Options {
	return Options{
		Namespace: "",
		Kind:      "",
		Name:      "",
		Version:   "",
		UID:       "",
		Group:     "",
		Weight:    0.0,
		Entity:    "",
		Attrs:     make(map[string]string),
	}
}
