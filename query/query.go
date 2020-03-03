package query

type Options struct {
	Namespace string
	Kind      string
	Name      string
}

type Option func(*Options)

func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

func Kind(kind string) Option {
	return func(o *Options) {
		o.Kind = kind
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func NewOptions() Options {
	return Options{
		Namespace: "",
		Kind:      "",
		Name:      "",
	}
}
