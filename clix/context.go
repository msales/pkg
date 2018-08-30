package clix

import (
	"context"

	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"github.com/urfave/cli"
)

type Ctx interface {
	Bool(name string) bool
	String(name string) string
	StringSlice(name string) []string
}

type CtxOptionFunc func(ctx *Context)

type ctxContext context.Context

type Context struct {
	*cli.Context
	ctxContext
}

func NewContext(c *cli.Context, opts ...CtxOptionFunc) (*Context, error) {
	ctx := &Context{
		Context:    c,
		ctxContext: context.Background(),
	}

	if len(opts) == 0 {
		l, err := NewLogger(ctx)
		if err != nil {
			return nil, err
		}

		s, err := NewStats(ctx, l)
		if err != nil {
			return nil, err
		}

		opts = append(opts, Logger(l), Stats(s))
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx, nil
}

func (c *Context) Close() error {
	s, ok := stats.FromContext(c)
	if ok {
		return s.Close()
	}

	return nil
}

func Logger(l log.Logger) CtxOptionFunc {
	return func(ctx *Context) {
		ctx.ctxContext = log.WithLogger(ctx.ctxContext, l)
	}
}

func Stats(s stats.Stats) CtxOptionFunc {
	return func(ctx *Context) {
		ctx.ctxContext = stats.WithStats(ctx.ctxContext, s)
	}
}
