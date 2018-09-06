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
		if -1 == strings.Index(r.URL.RawQuery, "[") {
			h.ServeHTTP(w, r)

			return
		}

		qs := r.URL.Query()
		normalizedQuery := make(url.Values, len(qs))

		for key, vals := range qs {
			newKey := normalizedQueryKey(key)

			for _, v := range vals {
				normalizedQuery.Add(newKey, v)
			}
		}

		r.URL.RawQuery = normalizedQuery.Encode()

		h.ServeHTTP(w, r)
	})
}

func normalizedQueryKey(key string) string {
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
