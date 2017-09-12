package stats

import (
	"context"
	"time"
)

// Timer represents a timer.
type Timer interface {
	// Start starts the timer.
	Start()
	// Done stops the timer and submits the Timing metric.
	Done()
}

type timer struct {
	start time.Time
	ctx   context.Context
	name  string
	rate  float32
	tags  map[string]string
}

// Time is a shorthand for Timing.
func Time(ctx context.Context, name string, rate float32, tags map[string]string) Timer {
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
