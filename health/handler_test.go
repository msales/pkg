package health_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v3/health"
	"github.com/stretchr/testify/assert"
)

func TestHandler_With(t *testing.T) {
	h := health.NewHandler()
	r1 := &testReporter{}
	r2 := &testReporter{}

	h.With(r1, r2)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	h.ServeHTTP(w, req)
	assert.True(t, r1.called)
	assert.True(t, r2.called)
}

func TestHandler_WithErrors(t *testing.T) {
	h := health.NewHandler()
	h.With(&testReporter{err: errors.New("test")})

	h.WithErrors()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, "test\n", w.Body.String())
}

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		reporterErrs []error
		code         int
	}{
		{[]error{nil}, http.StatusOK},
		{[]error{nil, nil}, http.StatusOK},
		{[]error{errors.New("")}, http.StatusServiceUnavailable},
		{[]error{errors.New(""), nil}, http.StatusServiceUnavailable},
		{[]error{nil, errors.New("")}, http.StatusServiceUnavailable},
		{[]error{errors.New(""), errors.New("")}, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		var reporters []health.Reporter

		for _, err := range tt.reporterErrs {
			r := &testReporter{err: err}
			reporters = append(reporters, r)
		}

		h := (health.NewHandler()).With(reporters...)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)

		h.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

type testReporter struct {
	called bool
	err    error
}

func (r *testReporter) IsHealthy() error {
	r.called = true

	return r.err
}
