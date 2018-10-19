package stats

import (
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/cactus/go-statsd-client/statsd/statsdtest"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsd(t *testing.T) {
	s, err := NewStatsd("127.0.0.1:1234", "test")
	assert.NoError(t, err)
	defer s.Close()

	assert.IsType(t, &Statsd{}, s)

	_, err = NewStatsd("127.0", "test")
	assert.Error(t, err)
}

func TestStatsd_Inc(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &Statsd{
		client: client,
	}

	s.Inc("test", 2, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestStatsd_Dec(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &Statsd{
		client: client,
	}

	s.Dec("test", 2, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "-2", sent[0].Value)
}

func TestStatsd_Gauge(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &Statsd{
		client: client,
	}

	s.Gauge("test", 2.0, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestStatsd_Timing(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &Statsd{
		client: client,
	}

	s.Timing("test", time.Second, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "1000", sent[0].Value)
}

func TestNewBufferedStatsd(t *testing.T) {
	s, err := NewBufferedStatsd("127.0.0.1:1234", "test", WithFlushInterval(time.Second), WithFlushBytes(1))
	assert.NoError(t, err)
	defer s.Close()

	assert.IsType(t, &BufferedStatsd{}, s)
	assert.Equal(t, time.Second, s.flushInterval)
	assert.Equal(t, 1, s.flushBytes)

	_, err = NewBufferedStatsd("127.0", "test")
	assert.Error(t, err)
}

func TestBufferedStatsd_Inc(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &BufferedStatsd{
		client: client,
	}

	s.Inc("test", 2, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestBufferedStatsd_Dec(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &BufferedStatsd{
		client: client,
	}

	s.Dec("test", 2, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "-2", sent[0].Value)
}

func TestBufferedStatsd_Gauge(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &BufferedStatsd{
		client: client,
	}

	s.Gauge("test", 2.0, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestBufferedStatsd_Timing(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &BufferedStatsd{
		client: client,
	}

	s.Timing("test", time.Second, 1.0, "test", "test")

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "1000", sent[0].Value)
}

func TestFormatStatsdTags(t *testing.T) {
	tags := []interface{}{
		"test", "test",
		"foo", "bar",
		"test", "baz",
	}

	assert.Equal(t, "", formatStatsdTags(nil))
	assert.Equal(t, "", formatStatsdTags([]interface{}{}))

	got := formatStatsdTags(tags)
	assert.NotContains(t, got, ",test=test")
	assert.Contains(t, got, ",test=baz")
	assert.Contains(t, got, ",foo=bar")
}

func TestFormatStatsdTags_Tags(t *testing.T) {
	tags := Tags{}.
		With("test", "test").
		With("foo", "bar").
		With("test", "test")

	got := formatStatsdTags([]interface{}{tags})
	assert.Contains(t, got, ",test=test")
	assert.Contains(t, got, ",foo=bar")
}

func TestFormatStatsdTags_Uneven(t *testing.T) {
	tags := []interface{}{
		"test", "test",
		"foo",
	}

	defer func() {
		assert.NotNil(t, recover())
	}()

	formatStatsdTags(tags)

	assert.Fail(t, "the test should have panicked on an uneven number of tags")
}

func BenchmarkFormatStatsdTags(b *testing.B) {
	tags := []interface{}{
		"string", "test",
		"float", 1.2,
		"int", 1,
		"bool", true,
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		formatStatsdTags(tags)
	}
}
