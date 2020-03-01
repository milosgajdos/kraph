package kraph

type Options struct {
	Namespace string
}

type Option func(*Options)

func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

type LinkOptions struct {
	Namespace bool
}

type LinkOption func(*LinkOptions)

func NamespaceLink(ns bool) LinkOption {
	return func(o *LinkOptions) {
		o.Namespace = ns
	}
}
