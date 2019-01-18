package middleware

import (
	"net/http"
	"net/url"
)

// StringTransformerFunc represents a function that transforms a string into another string.
type StringTransformerFunc func(string) string

type transformationRule struct {
	keys []string
	fn   StringTransformerFunc
}

// QueryTransformer transforms query parameters with the registered functions.
type QueryTransformer struct {
	rules []transformationRule
}

// Register registers a function as a transformer for the given set of keys.
func (t *QueryTransformer) Register(keys []string, fn StringTransformerFunc) {
	t.rules = append(t.rules, transformationRule{keys: keys, fn: fn})
}

// WithQueryTransformer returns an http handler that transforms query parameters using the passed transformer.
func WithQueryTransformer(h http.Handler, t QueryTransformer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()

		for _, rule := range t.rules {
			transformKeys(values, rule.keys, rule.fn)
		}

		req.URL.RawQuery = values.Encode()

		h.ServeHTTP(w, req)
	}
}

// WithQueryTransformerFunc returns an http handler that transforms query parameters using the passed function.
func WithQueryTransformerFunc(h http.Handler, keys []string, fn StringTransformerFunc) http.Handler {
	if len(keys) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()

		transformKeys(values, keys, fn)

		req.URL.RawQuery = values.Encode()

		h.ServeHTTP(w, req)
	})
}

func transformKeys(values url.Values, keys []string, fn StringTransformerFunc) {
	for _, key := range keys {
		vs := values[key]
		transformed := make([]string, len(vs))

		for j, v := range vs {
			transformed[j] = fn(v)
		}

		values[key] = transformed
	}
}