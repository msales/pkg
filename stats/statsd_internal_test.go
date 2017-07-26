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

	assert.IsType(t, &Statsd{}, s)
}

func TestStatsd_Inc(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &Statsd{
		client: client,
	}

	s.Inc("test", 2, 1.0, map[string]string{"test": "test"})

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

	s.Dec("test", 2, 1.0, map[string]string{"test": "test"})

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

	s.Gauge("test", 2.0, 1.0, map[string]string{"test": "test"})

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

	s.Timing("test", time.Second, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "1000", sent[0].Value)
}

func TestNewBufferedStatsd(t *testing.T) {
	s, err := NewBufferedStatsd("127.0.0.1:1234", "test", nil)
	assert.NoError(t, err)

	assert.IsType(t, &BufferedStatsd{}, s)
}

func TestBufferedStatsd_Inc(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	s := &BufferedStatsd{
		client: client,
	}

	s.Inc("test", 2, 1.0, map[string]string{"test": "test"})

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

	s.Dec("test", 2, 1.0, map[string]string{"test": "test"})

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

	s.Gauge("test", 2.0, 1.0, map[string]string{"test": "test"})

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

	s.Timing("test", time.Second, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "1000", sent[0].Value)
}

func TestFormatTags(t *testing.T) {
	tags := map[string]string{
		"test": "test",
		"foo":  "bar",
	}

	assert.Equal(t, "", formatTags(nil))
	assert.Equal(t, "", formatTags(map[string]string{}))

	got := formatTags(tags)
	assert.Contains(t, got, ",test=test")
	assert.Contains(t, got, ",foo=bar")
}
