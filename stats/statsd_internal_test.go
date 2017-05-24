package stats

import (
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/cactus/go-statsd-client/statsd/statsdtest"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsd(t *testing.T) {
	statsd, err := NewStatsd("127.0.0.1:1234", "test")
	assert.NoError(t, err)

	assert.IsType(t, &Statsd{}, statsd)
}

func TestStatsd_Inc(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	statsd := &Statsd{
		client: client,
	}

	statsd.Inc("test", 2, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestStatsd_Dec(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	statsd := &Statsd{
		client: client,
	}

	statsd.Dec("test", 2, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "-2", sent[0].Value)
}

func TestStatsd_Gauge(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	statsd := &Statsd{
		client: client,
	}

	statsd.Gauge("test", 2.0, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "2", sent[0].Value)
}

func TestStatsd_Timing(t *testing.T) {
	sender := statsdtest.NewRecordingSender()
	client, err := statsd.NewClientWithSender(sender, "test")
	assert.NoError(t, err)

	statsd := &Statsd{
		client: client,
	}

	statsd.Timing("test", time.Second, 1.0, map[string]string{"test": "test"})

	sent := sender.GetSent()
	assert.Len(t, sent, 1)
	assert.Equal(t, "test.test,test=test", sent[0].Stat)
	assert.Equal(t, "1000", sent[0].Value)
}

func TestStatsd_FormatTags(t *testing.T) {
	statsd := &Statsd{}
	tags := map[string]string{
		"test": "test",
		"foo":  "bar",
	}

	assert.Equal(t, "", statsd.formatTags(map[string]string{}))

	got := statsd.formatTags(tags)
	assert.Contains(t, got, ",test=test")
	assert.Contains(t, got, ",foo=bar")
}
