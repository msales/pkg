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
func (s *L2met) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	tags = deduplicateTags(normalizeTags(tags))
	ctx := append([]interface{}{
		s.formatL2metKey(name, "count"),
		fmt.Sprintf("%d", value),
	}, tags...)
	s.log.Info("", ctx...)

	return nil
}

// Dec decrements a count by the value.
func (s *L2met) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	tags = deduplicateTags(normalizeTags(tags))
	ctx := append([]interface{}{
		s.formatL2metKey(name, "count"),
		fmt.Sprintf("-%d", value),
	}, tags...)
	s.log.Info("", ctx...)

	return nil
}

// Gauge measures the value of a metric.
func (s *L2met) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	tags = deduplicateTags(normalizeTags(tags))
	ctx := append([]interface{}{
		s.formatL2metKey(name, "sample"),
		fmt.Sprintf("%v", value),
	}, tags...)
	s.log.Info("", ctx...)

	return nil
}

// Timing sends the value of a Duration.
func (s *L2met) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	tags = deduplicateTags(normalizeTags(tags))
	ctx := append([]interface{}{
		s.formatL2metKey(name, "measure"),
		formatDuration(value),
	}, tags...)
	s.log.Info("", ctx...)

	return nil
}

// Close closes the client and flushes buffered stats, if applicable
func (s *L2met) Close() error {
	return nil
}

// formatL2metKey creates an l2met compatible ctx key.
func (s *L2met) formatL2metKey(name, measure string) string {
	if s.prefix != "" {
		name = strings.Join([]string{s.prefix, name}, ".")
	}

	var buf bytes.Buffer
	buf.WriteString(measure)
	buf.WriteString("#")
	buf.WriteString(name)

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
