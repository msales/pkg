package stats

import (
	"context"
	"io"
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
	io.Closer

	// Inc increments a count by the value.
	Inc(name string, value int64, rate float32, tags ...interface{}) error

	// Dec decrements a count by the value.
	Dec(name string, value int64, rate float32, tags ...interface{}) error

	// Gauge measures the value of a metric.
	Gauge(name string, value float64, rate float32, tags ...interface{}) error

	// Timing sends the value of a Duration.
	Timing(name string, value time.Duration, rate float32, tags ...interface{}) error
}

// WithStats sets Stats in the context.
func WithStats(ctx context.Context, stats Stats) context.Context {
	if stats == nil {
		stats = Null
	}
	return context.WithValue(ctx, ctxKey, stats)
}

// FromContext returns the instance of Stats in the context.
func FromContext(ctx context.Context) (Stats, bool) {
	stats, ok := ctx.Value(ctxKey).(Stats)
	return stats, ok
}

// Inc increments a count by the value.
func Inc(ctx context.Context, name string, value int64, rate float32, tags ...interface{}) error {
	return withStats(ctx, func(s Stats) error {
		return s.Inc(name, value, rate, tags...)
	})
}

// Dec decrements a count by the value.
func Dec(ctx context.Context, name string, value int64, rate float32, tags ...interface{}) error {
	return withStats(ctx, func(s Stats) error {
		return s.Dec(name, value, rate, tags...)
	})
}

// Gauge measures the value of a metric.
func Gauge(ctx context.Context, name string, value float64, rate float32, tags ...interface{}) error {
	return withStats(ctx, func(s Stats) error {
		return s.Gauge(name, value, rate, tags...)
	})
}

// Timing sends the value of a Duration.
func Timing(ctx context.Context, name string, value time.Duration, rate float32, tags ...interface{}) error {
	return withStats(ctx, func(s Stats) error {
		return s.Timing(name, value, rate, tags...)
	})
}

// Close closes the client and flushes buffered stats, if applicable
func Close(ctx context.Context) error {
	return withStats(ctx, func(s Stats) error {
		return s.Close()
	})
}

func withStats(ctx context.Context, fn func(s Stats) error) error {
	if s, ok := FromContext(ctx); ok {
		return fn(s)
	}
	return fn(Null)
}

type nullStats struct{}

func (s nullStats) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	return nil
}

func (s nullStats) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	return nil
}

func (s nullStats) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	return nil
}

func (s nullStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	return nil
}

func (s nullStats) Close() error {
	return nil
}

// TaggedStats wraps a Stats instance applying tags to all metrics.
type TaggedStats struct {
	stats Stats
	tags  []interface{}
}

// NewTaggedStats creates a new TaggedStats instance.
func NewTaggedStats(stats Stats, tags ...interface{}) *TaggedStats {
	return &TaggedStats{
		stats: stats,
		tags:  normalizeTags(tags),
	}
}

// Inc increments a count by the value.
func (s TaggedStats) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	return s.stats.Inc(name, value, rate, mergeTags(tags, s.tags)...)
}

// Dec decrements a count by the value.
func (s TaggedStats) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	return s.stats.Dec(name, value, rate, mergeTags(tags, s.tags)...)
}

// Gauge measures the value of a metric.
func (s TaggedStats) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	return s.stats.Gauge(name, value, rate, mergeTags(tags, s.tags)...)
}

// Timing sends the value of a Duration.
func (s TaggedStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	return s.stats.Timing(name, value, rate, mergeTags(tags, s.tags)...)
}

// Close closes the client and flushes buffered stats, if applicable
func (s TaggedStats) Close() error {
	return s.stats.Close()
}

func normalizeTags(tags []interface{}) []interface{} {
	// If Tags object was passed, then expand it
	if len(tags) == 1 {
		if ctxMap, ok := tags[0].(Tags); ok {
			tags = ctxMap.toArray()
		}
	}

	// tags need to be even as they are key/value pairs
	if len(tags)%2 != 0 {
		panic("odd number of tags")
	}

	return tags
}

func mergeTags(prefix, suffix []interface{}) []interface{} {
	newTags := make([]interface{}, len(prefix)+len(suffix))

	n := copy(newTags, prefix)
	copy(newTags[n:], suffix)

	return newTags
}

func deduplicateTags(tags []interface{}) []interface{} {
	res := make([]interface{}, 0, len(tags))
Loop:
	for i := 0; i < len(tags); i += 2 {
		for j := 0; j < len(res); j += 2 {
			if tags[i] == res[j] {
				res[j+1] = tags[i+1]
				continue Loop
			}
		}

		res = append(res, tags[i])
		res = append(res, tags[i+1])
	}

	return res
}
