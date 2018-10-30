package stats_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/assert"
)

func TestPrometheus_Handler(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	h := s.Handler()

	assert.Implements(t, (*http.Handler)(nil), h)
}

func TestPrometheus_Inc(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	err := s.Inc("test", 2, 1.0, "test", "test")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	s.Handler().ServeHTTP(rr, req)

	assert.NoError(t, err)
	assert.Contains(t, rr.Body.String(), "test_test_test{test=\"test\"} 2")
}

func TestPrometheus_Dec(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	err := s.Dec("test", 2, 1.0, "test", "test")

	assert.Error(t, err)
}

func TestPrometheus_Gauge(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	err := s.Gauge("test", 2.1, 1.0, "test", "test")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	s.Handler().ServeHTTP(rr, req)

	assert.NoError(t, err)
	assert.Contains(t, rr.Body.String(), "test_test_test{test=\"test\"} 2.1")
}

func TestPrometheus_Timing(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	err := s.Timing("test", 1234500*time.Nanosecond, 1.0, "test", "test")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	s.Handler().ServeHTTP(rr, req)

	assert.NoError(t, err)
	assert.Contains(t, rr.Body.String(), "test_test_test{test=\"test\",quantile=\"0.5\"} 1.234")
	assert.Contains(t, rr.Body.String(), "test_test_test{test=\"test\",quantile=\"0.9\"} 1.234")
	assert.Contains(t, rr.Body.String(), "test_test_test{test=\"test\",quantile=\"0.99\"} 1.234")
	assert.Contains(t, rr.Body.String(), "test_test_test_sum{test=\"test\"} 1.234")
	assert.Contains(t, rr.Body.String(), "test_test_test_count{test=\"test\"} 1")
}

func TestPrometheus_Close(t *testing.T) {
	s := stats.NewPrometheus("test.test")

	err := s.Close()

	assert.NoError(t, err)
}
