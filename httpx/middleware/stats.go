package middleware

import (
	"net/http"
	"strconv"

	"github.com/msales/pkg/stats"
)

// WithRequestStats collects statistics about the request.
func WithRequestStats(h http.Handler, transformers ...PathTransformationFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for _, fn := range transformers {
			path = fn(path)
		}

		s, ok := stats.FromContext(r.Context())
		if !ok {
			s = stats.Null
		}

		s.Inc("request.start", 1, 1.0,
			"method", r.Method,
			"path", path,
		)

		rw := NewResponseWriter(w)
		h.ServeHTTP(rw, r)

		s.Inc("request.complete", 1, 1.0,
			"method", r.Method,
			"path", path,
			"status", strconv.Itoa(rw.Status()),
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

func ClearPath(string) string {
	return ""
}
