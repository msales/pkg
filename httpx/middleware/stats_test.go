package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"time"

	"github.com/msales/pkg/httpx/middleware"
	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestWithRequestStats(t *testing.T) {
	s := new(MockStats)
	s.On("Inc", "request.start", int64(1), float32(1.0), []interface{}{
		"method", "GET",
		"path",   "/",
	}).Return(nil).Once()
	s.On("Inc", "request.complete", int64(1), float32(1.0), []interface{}{
		"status", "0",
	}).Return(nil).Once()

	m := middleware.WithRequestStats(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	ctx := stats.WithStats(context.Background(), s)
	m.ServeHTTP(resp, req.WithContext(ctx))

	s.AssertExpectations(t)
}

func TestWithResponseTime(t *testing.T) {
	s := new(MockStats)
	s.On("Timing", "response.time", mock.Anything, float32(1.0), mock.Anything).Return(nil)

	m := middleware.WithResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	ctx := stats.WithStats(context.Background(), s)
	m.ServeHTTP(resp, req.WithContext(ctx))

	s.AssertExpectations(t)
}

type MockStats struct {
	mock.Mock
}

func (m *MockStats) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	args := m.Called(name, value, rate, tags)
	return args.Error(0)
}

func (m *MockStats) Close() error {
	args := m.Called()
	return args.Error(0)
}
