package k8s

type Options struct {
	Namespace string
}

type Option func(*Options)

func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}
