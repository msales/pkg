package stats_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/assert"
)

func TestL2met_Inc(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Inc("test", 2, 1.0, "foo", "bar")

	assert.Equal(t, "msg= count#test.test=2 foo=bar", l.Render())
}

func TestL2met_IncWithRate(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates(), stats.UseSampler(testSampler))

	s.Inc("test", 2, 0.1, "foo", "bar")

	assert.Equal(t, "msg= count#test.test@0.1=2 foo=bar", l.Render())
}

func TestL2met_Dec(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Dec("test", 2, 1.0, "foo", "bar")

	assert.Equal(t, "msg= count#test.test=-2 foo=bar", l.Render())
}

func TestL2met_DecWithRate(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates(), stats.UseSampler(testSampler))

	s.Dec("test", 2, 0.1, "foo", "bar")

	assert.Equal(t, "msg= count#test.test@0.1=-2 foo=bar", l.Render())
}

func TestL2met_Gauge(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Gauge("test", 2.1, 1.0, "foo", "bar")

	assert.Equal(t, "msg= sample#test.test=2.1 foo=bar", l.Render())
}

func TestL2met_GaugeWithRate(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates(), stats.UseSampler(testSampler))

	s.Gauge("test", 2.1, 0.1, "foo", "bar")

	assert.Equal(t, "msg= sample#test.test@0.1=2.1 foo=bar", l.Render())
}

func TestL2met_Timing(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 2*time.Second, 1.0, "foo", "bar")

	assert.Equal(t, "msg= measure#test.test=2000ms foo=bar", l.Render())
}

func TestL2met_TimingWithRate(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates(), stats.UseSampler(testSampler))

	s.Timing("test", 2*time.Second, 0.1, "foo", "bar")

	assert.Equal(t, "msg= measure#test.test@0.1=2000ms foo=bar", l.Render())
}

func TestL2met_Samples(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates(), stats.UseSampler(testNeverSampler))

	s.Timing("test", 2*time.Second, 0.1, "foo", "bar")

	assert.Equal(t, "msg=", l.Render())
}

func TestL2met_TimingFractions(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 1034500, 1.0, "foo", "bar")

	assert.Equal(t, "msg= measure#test.test=1.0345ms foo=bar", l.Render())
}

func TestL2met_TimingPartialFractions(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 1230*time.Microsecond, 1.0, "foo", "bar")

	assert.Equal(t, "msg= measure#test.test=1.23ms foo=bar", l.Render())
}

func TestL2met_Close(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	err := s.Close()

	assert.NoError(t, err)
}

func BenchmarkL2met_Inc(b *testing.B) {
	s := stats.NewL2met(log.Null, "test")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Inc("test", 2, 1.0, "foo", "bar")
	}
}

func BenchmarkL2met_Timing(b *testing.B) {
	s := stats.NewL2met(log.Null, "test")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Timing("test", 231, 1.0, "foo", "bar")
	}
}

func BenchmarkL2met_IncWithRate(b *testing.B) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test", stats.UseRates())

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Inc("test", 2, 0.1, "foo", "bar")
	}
}

func testSampler(f float32) bool {
	return true
}

func testNeverSampler(f float32) bool {
	return false
}

type testLogger struct {
	msg string
	ctx []interface{}
}

func (l *testLogger) Debug(msg string, ctx ...interface{}) {
	l.msg = msg
	l.ctx = ctx
}

func (l *testLogger) Info(msg string, ctx ...interface{}) {
	l.msg = msg
	l.ctx = ctx
}

func (l *testLogger) Error(msg string, ctx ...interface{}) {
	l.msg = msg
	l.ctx = ctx
}

func (l *testLogger) Render() string {
	var buf bytes.Buffer
	for i := 0; i < len(l.ctx); i += 2 {
		buf.WriteString(fmt.Sprintf("%v=%v ", l.ctx[i], l.ctx[i+1]))
	}

	return strings.Trim(fmt.Sprintf("msg=%s %s", l.msg, buf.String()), " ")
}
