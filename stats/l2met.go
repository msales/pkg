package stats

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/msales/pkg/v3/log"
)

// SamplerFunc represents a function that samples the L2met stats.
type SamplerFunc func(float32) bool

func defaultSampler(rate float32) bool {
	if rand.Float32() < rate {
		return true
	}
	return false
}

// L2metFunc represents a function that configures L2met.
type L2metFunc func(*L2met)

// UseRates turns on sample rates in l2met.
func UseRates() L2metFunc {
	return func(s *L2met) {
		s.useRates = true
	}
}

// UseSampler sets the sampler for l2met.
func UseSampler(sampler SamplerFunc) L2metFunc {
	return func(s *L2met) {
		s.sampler = sampler
	}
}

// L2met represents a l2met client.
type L2met struct {
	log    log.Logger
	prefix string

	useRates bool
	sampler  SamplerFunc
}

// NewL2met create a l2met instance.
func NewL2met(l log.Logger, prefix string, opts ...L2metFunc) Stats {
	s := &L2met{
		log:     l,
		prefix:  prefix,
		sampler: defaultSampler,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Inc increments a count by the value.
func (s *L2met) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	s.render(
		"count",
		name,
		strconv.FormatInt(value, 10),
		rate,
		tags,
	)

	return nil
}

// Dec decrements a count by the value.
func (s *L2met) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	s.render(
		"count",
		name,
		"-"+strconv.FormatInt(value, 10),
		rate,
		tags,
	)

	return nil
}

// Gauge measures the value of a metric.
func (s *L2met) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	s.render(
		"sample",
		name,
		strconv.FormatFloat(value, 'g', -1, 64),
		rate,
		tags,
	)

	return nil
}

// Timing sends the value of a Duration.
func (s *L2met) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	s.render(
		"measure",
		name,
		formatDuration(value),
		rate,
		tags,
	)

	return nil
}

// render outputs the metric to the logger
func (s *L2met) render(measure, name, value string, rate float32, tags []interface{}) {
	if !s.includeStat(rate) {
		return
	}

	tags = deduplicateTags(normalizeTags(tags))

	ctx := make([]interface{}, len(tags)+2)
	ctx[0] = s.formatL2metKey(name, measure) + s.formatL2metRate(rate)
	ctx[1] = value
	copy(ctx[2:], tags)

	s.log.Info("", ctx...)
}

func (s *L2met) includeStat(rate float32) bool {
	if !s.useRates || rate == 1.0 {
		return true
	}
	return s.sampler(rate)
}

// Close closes the client and flushes buffered stats, if applicable
func (s *L2met) Close() error {
	return nil
}

// formatL2metKey creates an l2met compatible ctx key.
func (s *L2met) formatL2metKey(name, measure string) string {
	if s.prefix != "" {
		return measure + "#" + s.prefix + "." + name
	}

	return measure + "#" + name
}

// formatL2metKey creates an l2met compatible rate suffix.
func (s *L2met) formatL2metRate(rate float32) string {
	if !s.useRates || rate == 1.0 {
		return ""
	}

	return "@" + strconv.FormatFloat(float64(rate), 'f', -1, 32)
}

// formatDuration converts duration into fractional milliseconds
// with no trailing zeros.
func formatDuration(d time.Duration) string {
	i := uint64(d / time.Millisecond)
	p := uint64(d % time.Millisecond / 1000)

	if p > 0 {
		for {
			if p%10 == 0 {
				p /= 10
				continue
			}
			break
		}

		return strconv.FormatUint(i, 10) + "." + strconv.FormatUint(p, 10) + "ms"
	}

	return strconv.FormatUint(i, 10) + "ms"
}
