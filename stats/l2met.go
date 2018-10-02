package stats

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/msales/pkg/v2/log"
)

// L2met represents a l2met client.
type L2met struct {
	log    log.Logger
	prefix string
}

// NewL2met create a l2met instance.
func NewL2met(l log.Logger, prefix string) Stats {
	return &L2met{
		log:    l,
		prefix: prefix,
	}
}

// Inc increments a count by the value.
func (s L2met) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	msg := s.formatL2metMetric(name, fmt.Sprintf("%d", value), "count", tags)
	s.log.Info(msg)

	return nil
}

// Dec decrements a count by the value.
func (s L2met) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	msg := s.formatL2metMetric(name, fmt.Sprintf("-%d", value), "count", tags)
	s.log.Info(msg)

	return nil
}

// Gauge measures the value of a metric.
func (s L2met) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	msg := s.formatL2metMetric(name, fmt.Sprintf("%v", value), "sample", tags)
	s.log.Info(msg)

	return nil
}

// Timing sends the value of a Duration.
func (s L2met) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	msg := s.formatL2metMetric(name, formatDuration(value), "measure", tags)
	s.log.Info(msg)

	return nil
}

// Close closes the client and flushes buffered stats, if applicable
func (s L2met) Close() error {
	return nil
}

func (s L2met) formatL2metMetric(name, value, measure string, tags []interface{}) string {
	if s.prefix != "" {
		name = strings.Join([]string{s.prefix, name}, ".")
	}

	var buf bytes.Buffer
	buf.WriteString(formatL2metTags(tags))
	buf.WriteString(measure)
	buf.WriteString("#")
	buf.WriteString(name)
	buf.WriteString("=")
	buf.WriteString(value)

	return buf.String()
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

		return fmt.Sprintf("%d.%dms", i, p)
	}

	return fmt.Sprintf("%dms", i)
}

// formatStatsdTags formats into an InfluxDB style string
func formatL2metTags(tags []interface{}) string {
	if len(tags) == 0 {
		return ""
	}

	tags = deduplicateTags(normalizeTags(tags))

	var buf bytes.Buffer
	for i := 0; i < len(tags); i += 2 {
		buf.WriteString(fmt.Sprintf("%v=%v ", tags[i], tags[i+1]))
	}

	return buf.String()
}
