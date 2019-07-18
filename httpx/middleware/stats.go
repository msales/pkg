package middleware

import (
	"net/http"
	"strconv"

	"github.com/msales/pkg/v3/stats"
)

// TagsFunc returns a set of tags from a request
type TagsFunc func(*http.Request) []interface{}

// DefaultTags extracts the method and path from the request.
func DefaultTags(r *http.Request) []interface{} {
	return []interface{}{
		"method", r.Method,
		"path", r.URL.Path,
	}
}

// WithRequestStats collects statistics about the request.
func WithRequestStats(h http.Handler, fns ...TagsFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := prepareTags(r, fns)

		s, ok := stats.FromContext(r.Context())
		if !ok {
			s = stats.Null
		}

		_ = s.Inc("request.start", 1, 1.0, tags...)

		rw := NewResponseWriter(w)
		h.ServeHTTP(rw, r)

		cpltTags := make([]interface{}, len(tags)+2)
		cpltTags[0] = "status"
		cpltTags[1] = strconv.FormatInt(int64(rw.Status()), 10)
		copy(cpltTags[2:], tags)
		_ = s.Inc("request.complete", 1, 1.0, cpltTags...)
	})
}

// WithResponseTime reports the response time.
func WithResponseTime(h http.Handler, fns ...TagsFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := prepareTags(r, fns)
		t := stats.Time(r.Context(), "response.time", 1.0, tags...)
		defer t.Done()

		h.ServeHTTP(w, r)
	})
}

// prepareTags resolves tags in accordance to provided functions and falls back to defaults in no custom tag functions were provided.
func prepareTags(r *http.Request, fns []TagsFunc) []interface{} {
	if len(fns) == 0 {
		return nil
	}

	var tags []interface{}
	for _, fn := range fns {
		tags = append(tags, fn(r)...)
	}

	return tags
}
