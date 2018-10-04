package stats_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestL2met_Inc(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Inc("test", 2, 1.0, "test", "test")

	assert.Equal(t, "msg= count#test.test=2 test=test", l.Render())
}

func TestL2met_Dec(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Dec("test", 2, 1.0, "test", "test")

	assert.Equal(t, "msg= count#test.test=-2 test=test", l.Render())
}

func TestL2met_Gauge(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Gauge("test", 2.1, 1.0, "test", "test")

	assert.Equal(t, "msg= sample#test.test=2.1 test=test", l.Render())
}

func TestL2met_Timing(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 2*time.Second+time.Nanosecond, 1.0, "test", "test")

	assert.Equal(t, "msg= measure#test.test=2000ms test=test", l.Render())
}

func TestL2met_TimingFractions(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 1234500*time.Nanosecond, 1.0, "test", "test")

	assert.Equal(t, "msg= measure#test.test=1.234ms test=test", l.Render())
}

func TestL2met_TimingPartialFractions(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 1230*time.Microsecond, 1.0, "test", "test")

	assert.Equal(t, "msg= measure#test.test=1.23ms test=test", l.Render())
}

func TestL2met_Close(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	err := s.Close()

	assert.NoError(t, err)
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