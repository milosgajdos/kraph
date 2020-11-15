package star

// Options provides gh options
type Options struct {
	Paging  int
	User    string
	Workers int
}

// Option is gh option
type Option func(*Options)

// Paging configures paging
func Paging(p int) Option {
	return func(o *Options) {
		o.Paging = p
	}
}

// User configures gh user
func User(u string) Option {
	return func(o *Options) {
		o.User = u
	}
}

// Workers configures paging
func Workers(w int) Option {
	return func(o *Options) {
		o.Workers = w
	}
}
