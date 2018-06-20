package stats_test

import (
	"testing"
	"time"

	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestL2met_Inc(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Inc("test", 2, 1.0, "test", "test")

	assert.Equal(t, "test=test count#test.test=2", l.msg)
}

func TestL2met_Dec(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Dec("test", 2, 1.0, "test", "test")

	assert.Equal(t, "test=test count#test.test=-2", l.msg)
}

func TestL2met_Gauge(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Gauge("test", 2.1, 1.0, "test", "test")

	assert.Equal(t, "test=test sample#test.test=2.1", l.msg)
}

func TestL2met_Timing(t *testing.T) {
	l := &testLogger{}
	s := stats.NewL2met(l, "test")

	s.Timing("test", 2*time.Second+2*time.Microsecond, 1.0, "test", "test")

	assert.Equal(t, "test=test measure#test.test=2000ms", l.msg)
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
