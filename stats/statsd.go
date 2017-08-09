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
func NewStatsd(addr, prefix string) (Stats, error) {
	c, err := statsd.NewClient(addr, prefix)
	if err != nil {
		return nil, err
	}

	return &Statsd{
		client: c,
	}, nil
}

// Inc increments a count by the value.
func (s Statsd) Inc(name string, value int64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Inc(name, value, rate)
}

// Dec decrements a count by the value.
func (s Statsd) Dec(name string, value int64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Dec(name, value, rate)
}

// Gauge measures the value of a metric.
func (s Statsd) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Gauge(name, int64(value), rate)
}

// Timing sends the value of a Duration.
func (s Statsd) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.TimingDuration(name, value, rate)
}

// BufferedStatsdFunc represents an configuration function for BufferedStatsd.
type BufferedStatsdFunc func(*BufferedStatsd)

// WithFlushInterval sets the maximum flushInterval for packet sending.
// Defaults to 300ms.
func WithFlushInterval(interval time.Duration) BufferedStatsdFunc {
	return func(s *BufferedStatsd) {
		s.flushInterval = interval
	}
}

// WithFlushBytes sets the maximum udp packet size that will be sent.
// Defaults to 1432 flushBytes.
func WithFlushBytes(bytes int) BufferedStatsdFunc {
	return func(s *BufferedStatsd) {
		s.flushBytes = bytes
	}
}

// BufferedStatsd represents a buffered statsd client.
type BufferedStatsd struct {
	client statsd.Statter

	flushInterval time.Duration
	flushBytes    int
}

// NewBufferedStatsd create a buffered Statsd instance.
func NewBufferedStatsd(addr, prefix string, opts ...BufferedStatsdFunc) (*BufferedStatsd, error) {
	s := &BufferedStatsd{}

	for _, o := range opts {
		o(s)
	}

	c, err := statsd.NewBufferedClient(addr, prefix, s.flushInterval, s.flushBytes)
	if err != nil {
		return nil, err
	}
	s.client = c

	return s, nil
}

// Inc increments a count by the value.
func (s BufferedStatsd) Inc(name string, value int64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Inc(name, value, rate)
}

// Dec decrements a count by the value.
func (s BufferedStatsd) Dec(name string, value int64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Dec(name, value, rate)
}

// Gauge measures the value of a metric.
func (s BufferedStatsd) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.Gauge(name, int64(value), rate)
}

// Timing sends the value of a Duration.
func (s BufferedStatsd) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	name += formatTags(tags)
	return s.client.TimingDuration(name, value, rate)
}

// formatTags formats into an InfluxDB style string
func formatTags(tags map[string]string) string {
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
