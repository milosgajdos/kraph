package k8s

// Options provides k8so options
type Options struct {
	Namespace string
}

// Option is k8s option
type Option func(*Options)

// Namespace configures namespace
func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}
