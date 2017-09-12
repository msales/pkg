package stats

import (
	"context"
	"time"
)

type key int

const (
	ctxKey key = iota
)

var (
	// Null is the null Stats instance.
	Null = &nullStats{}
)

// Stats represents a stats instance.
type Stats interface {
	// Inc increments a count by the value.
	Inc(name string, value int64, rate float32, tags map[string]string) error

	// Dec decrements a count by the value.
	Dec(name string, value int64, rate float32, tags map[string]string) error

	// Gauge measures the value of a metric.
	Gauge(name string, value float64, rate float32, tags map[string]string) error

	// Timing sends the value of a Duration.
	Timing(name string, value time.Duration, rate float32, tags map[string]string) error
}

// WithStats sets Stats in the context.
func WithStats(ctx context.Context, stats Stats) context.Context {
	return context.WithValue(ctx, ctxKey, stats)
}

// FromContext returns the instance of Stats in the context.
func FromContext(ctx context.Context) (Stats, bool) {
	stats, ok := ctx.Value(ctxKey).(Stats)
	return stats, ok
}

// Inc increments a count by the value.
func Inc(ctx context.Context, name string, value int64, rate float32, tags map[string]string) error {
	return withStats(ctx, func(s Stats) error {
		return s.Inc(name, value, rate, tags)
	})
}

// Dec decrements a count by the value.
func Dec(ctx context.Context, name string, value int64, rate float32, tags map[string]string) error {
	return withStats(ctx, func(s Stats) error {
		return s.Dec(name, value, rate, tags)
	})
}

// Gauge measures the value of a metric.
func Gauge(ctx context.Context, name string, value float64, rate float32, tags map[string]string) error {
	return withStats(ctx, func(s Stats) error {
		return s.Gauge(name, value, rate, tags)
	})
}

// Timing sends the value of a Duration.
func Timing(ctx context.Context, name string, value time.Duration, rate float32, tags map[string]string) error {
	return withStats(ctx, func(s Stats) error {
		return s.Timing(name, value, rate, tags)
	})
}

func withStats(ctx context.Context, fn func(s Stats) error) error {
	if s, ok := FromContext(ctx); ok {
		return fn(s)
	}
	return fn(Null)
}

type nullStats struct{}

func (s nullStats) Inc(name string, value int64, rate float32, tags map[string]string) error {
	return nil
}

func (s nullStats) Dec(name string, value int64, rate float32, tags map[string]string) error {
	return nil
}

func (s nullStats) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	return nil
}

func (s nullStats) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	return nil
}

// TaggedStats wraps a Stats instance applying tags to all metrics.
type TaggedStats struct {
	stats Stats
	tags  map[string]string
}

// NewTaggedStats creates a new TaggedStats instance.
func NewTaggedStats(stats Stats, tags map[string]string) *TaggedStats {
	return &TaggedStats{
		stats: stats,
		tags:  tags,
	}
}

// Inc increments a count by the value.
func (s TaggedStats) Inc(name string, value int64, rate float32, tags map[string]string) error {
	return s.stats.Inc(name, value, rate, s.collectTags(tags))
}

// Dec decrements a count by the value.
func (s TaggedStats) Dec(name string, value int64, rate float32, tags map[string]string) error {
	return s.stats.Dec(name, value, rate, s.collectTags(tags))
}

// Gauge measures the value of a metric.
func (s TaggedStats) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	return s.stats.Gauge(name, value, rate, s.collectTags(tags))
}

// Timing sends the value of a Duration.
func (s TaggedStats) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	return s.stats.Timing(name, value, rate, s.collectTags(tags))
}

func (s TaggedStats) collectTags(tags map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range s.tags {
		res[k] = v
	}
	for k, v := range tags {
		res[k] = v
	}

	return res
}
