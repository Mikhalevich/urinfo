package request

type options struct {
	ForceHTTP11 bool
}

type Option func(opts *options)

func WithForceHTTP11(force bool) Option {
	return func(opts *options) {
		opts.ForceHTTP11 = force
	}
}
