package owner

// Options provides k8s options
type Options struct {
	Namespace string
}

// Option is k8s option
type Option func(*Options)

// Namespace configures k8s namespace
func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}
