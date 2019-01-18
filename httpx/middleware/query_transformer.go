package middleware

import (
	"net/http"
	"net/url"
)

type StringTransformerFunc func(string) string

type transformationRule struct {
	keys []string
	fn   StringTransformerFunc
}

type QueryTransformer struct {
	rules []transformationRule
}

func (t *QueryTransformer) Register(keys []string, fn StringTransformerFunc) {
	t.rules = append(t.rules, transformationRule{keys: keys, fn: fn})
}

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