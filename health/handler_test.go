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
	h := &health.Handler{}
	r1 := &ReporterMock{}
	r2 := &ReporterMock{}

	h.With(r1, r2)

	assert.Contains(t, h.Reporters, r1)
	assert.Contains(t, h.Reporters, r2)
}

func TestHandler_WithErrors(t *testing.T) {
	h := &health.Handler{}

	assert.False(t, h.ShowErr)

	h.WithErrors()

	assert.True(t, h.ShowErr)
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
			r := ReporterMock{err}
			reporters = append(reporters, r)
		}

		h := (&health.Handler{}).With(reporters...)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)

		h.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

type ReporterMock struct {
	healthy error
}

func (r ReporterMock) IsHealthy() error {
	return r.healthy
}
