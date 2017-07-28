package stats_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/msales/pkg/stats"
)

func TestTaggedStats_Inc(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "test", int64(1), float32(1), map[string]string{"foo": "bar", "global": "foobar"}).Return(nil)

	s := stats.NewTaggedStats(m, map[string]string{"global": "foobar"})
	s.Inc("test", 1, 1.0, map[string]string{"foo": "bar"})

	m.AssertExpectations(t)
}

func TestTaggedStats_Dec(t *testing.T) {
	m := new(MockStats)
	m.On("Dec", "test", int64(1), float32(1), map[string]string{"foo": "bar", "global": "foobar"}).Return(nil)

	s := stats.NewTaggedStats(m, map[string]string{"global": "foobar"})
	s.Dec("test", 1, 1.0, map[string]string{"foo": "bar"})

	m.AssertExpectations(t)
}

func TestTaggedStats_Gauge(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "test", float64(1), float32(1), map[string]string{"foo": "bar", "global": "foobar"}).Return(nil)

	s := stats.NewTaggedStats(m, map[string]string{"global": "foobar"})
	s.Gauge("test", 1.0, 1.0, map[string]string{"foo": "bar"})

	m.AssertExpectations(t)
}

func TestTaggedStats_Timing(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", time.Millisecond, float32(1), map[string]string{"foo": "bar", "global": "foobar"}).Return(nil)

	s := stats.NewTaggedStats(m, map[string]string{"global": "foobar"})
	s.Timing("test", time.Millisecond, 1.0, map[string]string{"foo": "bar"})

	m.AssertExpectations(t)
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags map[string]string) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Dec(name string, value int64, rate float32, tags map[string]string) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags map[string]string) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags map[string]string) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}
