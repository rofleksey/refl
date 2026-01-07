package eval

type Options struct {
	disableEvents bool
	disableEval   bool
}

type Option interface {
	Apply(*Options)
}

type OptionDisableEvents struct{}

func (OptionDisableEvents) Apply(opts *Options) {
	opts.disableEvents = true
}

type OptionDisableEval struct{}

func (OptionDisableEval) Apply(opts *Options) {
	opts.disableEval = true
}

type OptionSetOptions struct {
	opts Options
}

func (o OptionSetOptions) Apply(opts *Options) {
	*opts = o.opts
}
