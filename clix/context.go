package clix

import (
	"context"

	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"gopkg.in/urfave/cli.v1"
)

type ctxContext context.Context

// Context represents an application context.
type Context struct {
	*cli.Context
	ctxContext
}

// NewContext creates a new Context from the CLI Context.
func NewContext(c *cli.Context, opts ...ContextFunc) (*Context, error) {
	ctx := &Context{
		Context:    c,
		ctxContext: context.Background(),
	}

	for _, opt := range opts {
		opt(ctx)
	}

	if _, ok := log.FromContext(ctx); !ok {
		l, err := NewLogger(ctx.Context)
		if err != nil {
			return nil, err
		}

		WithLogger(l)(ctx)
	}

	if _, ok := stats.FromContext(ctx); !ok {
		l, _ := log.FromContext(ctx) // guaranteed to have a WithLogger instance here
		s, err := NewStats(ctx.Context, l)
		if err != nil {
			return nil, err
		}

		WithStats(s)(ctx)
	}

	return ctx, nil
}

// Close closes the context.
func (c *Context) Close() error {
	s, ok := stats.FromContext(c)
	if ok {
		return s.Close()
	}

	return nil
}

// ContextFunc configures the Context.
type ContextFunc func(ctx *Context)

// WithLogger sets the logger instance on the Context.
func WithLogger(l log.Logger) ContextFunc {
	return func(ctx *Context) {
		ctx.ctxContext = log.WithLogger(ctx.ctxContext, l)
	}
}

// WithStats set the stats instance on the Context.
func WithStats(s stats.Stats) ContextFunc {
	return func(ctx *Context) {
		ctx.ctxContext = stats.WithStats(ctx.ctxContext, s)
	}
}
