package clix

import (
	"context"
	"io"

	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"gopkg.in/urfave/cli.v1"
)

// ContextFunc configures the Context.
type ContextFunc func(ctx *Context)

// WithLogger sets the logger instance on the Context.
func WithLogger(l log.Logger) ContextFunc {
	return func(ctx *Context) {
		ctx.logger = l
	}
}

// WithStats set the stats instance on the Context.
func WithStats(s stats.Stats) ContextFunc {
	return func(ctx *Context) {
		ctx.stats = s
	}
}

type ctxContext context.Context

// Context represents an application context.
type Context struct {
	*cli.Context
	ctxContext

	logger log.Logger
	stats  stats.Stats
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

	if ctx.logger == nil {
		l, err := NewLogger(ctx.Context)
		if err != nil {
			return nil, err
		}

		ctx.logger = l
	}

	if ctx.stats == nil {
		s, err := NewStats(ctx.Context, ctx.logger) // guaranteed to have a logger instance here
		if err != nil {
			return nil, err
		}

		WithStats(s)(ctx)
	}

	ctx.ctxContext = log.WithLogger(ctx.ctxContext, ctx.logger)
	ctx.ctxContext = stats.WithStats(ctx.ctxContext, ctx.stats)

	return ctx, nil
}

// Close closes the context.
func (c *Context) Close() error {
	if err := c.stats.Close(); err != nil {
		return err
	}

	if l, ok := c.logger.(io.Closer); ok {
		return l.Close()
	}

	return nil
}
