package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msales/pkg/v4/httpx/middleware"
	"github.com/msales/pkg/v4/stats"
	"github.com/stretchr/testify/mock"
)

func TestWithRequestStats(t *testing.T) {
	tests := []struct {
		path         string
		tagFuncs     []middleware.TagsFunc
		expectedTags []interface{}
	}{
		{"/test", nil, []interface{}{"method", "GET", "path", "/test"}},
		{"", nil, []interface{}{"method", "GET", "path", ""}},
		{"/test", []middleware.TagsFunc{testTags}, []interface{}{"method", "GET"}},
		{"", []middleware.TagsFunc{testTags}, []interface{}{"method", "GET"}},
	}

	for _, tt := range tests {
		s := new(MockStats)
		s.On("Inc", "request.start", int64(1), float32(1.0), tt.expectedTags).Return(nil).Once()
		s.On("Inc", "request.complete", int64(1), float32(1.0), append([]interface{}{
			"status", "0",
		}, tt.expectedTags...)).Return(nil).Once()

		m := middleware.WithRequestStats(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), tt.tagFuncs...)

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", tt.path, nil)

		ctx := stats.WithStats(context.Background(), s)
		m.ServeHTTP(resp, req.WithContext(ctx))

		s.AssertExpectations(t)
	}
}

func TestWithRequestStats_NoStats(t *testing.T) {
	m := middleware.WithRequestStats(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	m.ServeHTTP(resp, req.WithContext(context.Background()))
}

func TestWithResponseTime(t *testing.T) {
	tests := []struct {
		tagFuncs     []middleware.TagsFunc
		expectedTags []interface{}
	}{
		{
			tagFuncs:     nil,
			expectedTags: []interface{}{"method", "GET", "path", "/"},
		},
		{
			tagFuncs:     []middleware.TagsFunc{testTags},
			expectedTags: []interface{}{"method", "GET"},
		},
	}

	for _, tt := range tests {
		s := new(MockStats)
		s.On("Timing", "response.time", mock.Anything, float32(1.0), tt.expectedTags).Return(nil).Once()

		m := middleware.WithResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), tt.tagFuncs...)

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		ctx := stats.WithStats(context.Background(), s)
		m.ServeHTTP(resp, req.WithContext(ctx))

		s.AssertExpectations(t)
	}
}

func testTags(r *http.Request) []interface{} {
	return []interface{}{
		"method", r.Method,
	}
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
