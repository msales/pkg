package stats

import (
	"context"
	"time"
)

type timer struct {
	start time.Time
	ctx   context.Context
	name  string
	rate  float32
	tags  map[string]string
}

// Time is a shorthand for Timing.
func Time(ctx context.Context, name string, rate float32, tags map[string]string) *timer {
	t := &timer{ctx: ctx, name: name, rate: rate, tags: tags}
	t.Start()
	return t
}

func (t *timer) Start() {
	t.start = time.Now()
}

func (t *timer) Done() {
	Timing(t.ctx, t.name, time.Since(t.start), t.rate, t.tags)
}
