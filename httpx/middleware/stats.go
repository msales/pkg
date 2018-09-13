package middleware

import (
	"fmt"
	"net/http"

	"github.com/msales/pkg/stats"
)

// WithRequestStats collects statistics about the request.
func WithRequestStats(h http.Handler, transformers ...PathTransformationFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		for _, fn := range transformers {
			path = fn(path)
		}

		stats.Inc(r.Context(), "request.start", 1, 1.0,
			"method", r.Method,
			"path", path,
		)

		rw := NewResponseWriter(w)
		h.ServeHTTP(rw, r)

		stats.Inc(r.Context(), "request.complete", 1, 1.0,
			"status", fmt.Sprintf("%d", rw.Status()),
			"path", path,
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

type PathTransformationFunc func(path string) string

func Clear(string) string {
	return ""
}
