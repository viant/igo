package option

import (
	"github.com/viant/igo/notify"
)

// Options represents planner options
type Options struct {
	notify.StmtListener
	Tracker *notify.Tracker
	notify.ExprListener
	options []Option
}

// NewOptions returns a new options
func NewOptions(opts ...Option) *Options {
	o := &Options{
		options: opts,
	}
	o.Apply(opts)
	return o
}

func (o *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *Options) Merge(parent *Options, opts []Option) *Options {
	if parent == nil {
		o.Apply(opts)
		return o
	}
	o.StmtListener = parent.StmtListener
	o.Tracker = parent.Tracker
	o.ExprListener = parent.ExprListener
	o.Apply(opts)
	return o
}

// Option represents a planner option
type Option func(o *Options)

// WithStmtListener sets the statement listener
func WithStmtListener(listener notify.StmtListener) Option {
	return func(o *Options) {
		o.StmtListener = listener
	}
}

// WithTracker sets the tracker
func WithTracker(tracker *notify.Tracker) Option {
	return func(o *Options) {
		o.Tracker = tracker
	}
}

// WithExprListener returns an option for setting the expression listener
func WithExprListener(listener notify.ExprListener) Option {
	return func(o *Options) {
		o.ExprListener = listener
	}
}
