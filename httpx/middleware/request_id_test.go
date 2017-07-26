package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/httpx"
	"github.com/msales/pkg/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		header http.Header
		id     string
	}{
		{http.Header{http.CanonicalHeaderKey("X-Request-ID"): []string{"1234"}}, "1234"},
		{http.Header{http.CanonicalHeaderKey("Request-ID"): []string{"1234"}}, "1234"},
		{http.Header{http.CanonicalHeaderKey("Foo"): []string{"1234"}}, ""},
	}

	for _, tt := range tests {
		m := middleware.ExtractRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := httpx.RequestID(r.Context())

			assert.Equal(t, tt.id, requestID)
		}))

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = tt.header

		m.ServeHTTP(resp, req)
	}
}
