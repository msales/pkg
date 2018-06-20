package middleware

import (
	"fmt"
	"net/http"

	"github.com/msales/pkg/stats"
)

// WithRequestStats collects statistics about the request.
func WithRequestStats(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats.Inc(r.Context(), "request.start", 1, 1.0,
			"method", r.Method,
			"path", r.URL.Path,
		)

		rw := NewResponseWriter(w)
		h.ServeHTTP(rw, r)

		stats.Inc(r.Context(), "request.complete", 1, 1.0,
			"status", fmt.Sprintf("%d", rw.Status()),
		)
	})
}

// WithResponseTime reports the response time.
func WithResponseTime(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := stats.Time(r.Context(), "response.time", 1.0)
		defer t.Done()

		h.ServeHTTP(w, r)
	})
}
