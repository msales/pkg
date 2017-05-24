package stats

import (
	"bytes"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
)

// Statsd represents a statsd client.
type Statsd struct {
	client statsd.Statter
}

// NewStatsd create a Statsd instance.
func NewStatsd(addr, prefix string) (*Statsd, error) {
	c, err := statsd.New(addr, prefix)
	if err != nil {
		return nil, err
	}

	return &Statsd{
		client: c,
	}, nil
}

// Inc increments a count by the value.
func (s Statsd) Inc(name string, value int64, rate float32, tags map[string]string) error {
	name += s.formatTags(tags)
	return s.client.Inc(name, value, rate)
}

// Dec decrements a count by the value.
func (s Statsd) Dec(name string, value int64, rate float32, tags map[string]string) error {
	name += s.formatTags(tags)
	return s.client.Dec(name, value, rate)
}

// Gauge measures the value of a metric.
func (s Statsd) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	name += s.formatTags(tags)
	return s.client.Gauge(name, int64(value), rate)
}

// Timing sends the value of a Duration.
func (s Statsd) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	name += s.formatTags(tags)
	return s.client.TimingDuration(name, value, rate)
}

// formatTags formats into an InfluxDB style string
func (s Statsd) formatTags(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for k, v := range tags {
		buf.WriteString(",")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}

	return buf.String()
}
