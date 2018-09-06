package middleware

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// WithQueryNormalizer fixes wrong php query string handling as array
func WithQueryNormalizer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if i := strings.Index(r.URL.RawQuery, "["); i == -1 {
			h.ServeHTTP(w, r)

			return
		}

		q := r.URL.Query()
		normalizedQuery := make(url.Values, 0)

		for key, qVal := range q {
			newKey := getNormalizedValue(key)

			for _, v := range qVal {
				normalizedQuery.Add(newKey, v)
			}
		}

		r.URL.RawQuery = normalizedQuery.Encode()

		h.ServeHTTP(w, r)
	})
}

func getNormalizedValue(key string) string {
	strs := strings.Split(key, "[")
	if len(strs) == 0 {
		return key
	}

	for i, str := range strs {
		if _, err := strconv.Atoi(str[:len(str)-1]); err == nil {
			strs[i] = "]"
		}
	}
	return strings.Join(strs, "[")
}
