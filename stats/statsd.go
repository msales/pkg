package stats

import (
	"fmt"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/msales/pkg/v3/bytes"
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
func (s *Statsd) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Inc(name, value, rate)
}

// Dec decrements a count by the value.
func (s *Statsd) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Dec(name, value, rate)
}

// Gauge measures the value of a metric.
func (s *Statsd) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Gauge(name, int64(value), rate)
}

// Timing sends the value of a Duration.
func (s *Statsd) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.TimingDuration(name, value, rate)
}

// Close closes the client and flushes buffered stats, if applicable
func (s *Statsd) Close() error {
	return s.client.Close()
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
func (s *BufferedStatsd) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Inc(name, value, rate)
}

// Dec decrements a count by the value.
func (s *BufferedStatsd) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Dec(name, value, rate)
}

// Gauge measures the value of a metric.
func (s *BufferedStatsd) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.Gauge(name, int64(value), rate)
}

// Timing sends the value of a Duration.
func (s *BufferedStatsd) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	name += formatStatsdTags(tags)
	return s.client.TimingDuration(name, value, rate)
}

// Close closes the client and flushes buffered stats, if applicable
func (s *BufferedStatsd) Close() error {
	return s.client.Close()
}

var statsdPool = bytes.NewPool(100)

// formatStatsdTags formats into an InfluxDB style string
func formatStatsdTags(tags []interface{}) string {
	if len(tags) == 0 {
		return ""
	}

	tags = deduplicateTags(normalizeTags(tags))

	buf := statsdPool.Get()
	for i := 0; i < len(tags); i += 2 {
		buf.WriteByte(',')
		formatStatsdValue(buf, tags[i])
		buf.WriteByte('=')
		formatStatsdValue(buf, tags[i+1])
	}

	s := string(buf.Bytes())
	statsdPool.Put(buf)
	return s
}

// formatStatsdValue formats a value, adding it to the Buffer.
func formatStatsdValue(buf *bytes.Buffer, value interface{}) {
	if value == nil {
		buf.WriteString("nil")
		return
	}

	switch v := value.(type) {
	case bool:
		buf.AppendBool(v)
	case float32:
		buf.AppendFloat(float64(v), 'g', -1, 64)
	case float64:
		buf.AppendFloat(v, 'g', -1, 64)
	case int:
		buf.AppendInt(int64(v))
	case int8:
		buf.AppendInt(int64(v))
	case int16:
		buf.AppendInt(int64(v))
	case int32:
		buf.AppendInt(int64(v))
	case int64:
		buf.AppendInt(v)
	case uint:
		buf.AppendUint(uint64(v))
	case uint8:
		buf.AppendUint(uint64(v))
	case uint16:
		buf.AppendUint(uint64(v))
	case uint32:
		buf.AppendUint(uint64(v))
	case uint64:
		buf.AppendUint(v)
	case string:
		buf.WriteString(v)
	default:
		buf.WriteString(fmt.Sprintf("%+v", value))
	}
}
