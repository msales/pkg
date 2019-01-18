package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/v3/httpx/middleware"
	"github.com/stretchr/testify/assert"
)

func TestWithQueryTransformer(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	seenVals := make([]string, 0)
	fn := func(s string) string {
		seenVals = append(seenVals, s)
		return s
	}

	var tr middleware.QueryTransformer
	tr.Register([]string{"a", "b"}, fn)
	tr.Register([]string{"c"}, fn)
	tr.Register([]string{}, fn)

	h = middleware.WithQueryTransformer(h, tr)

	req, _ := http.NewRequest("GET", "/?a=va&b=vb&c=vc&d=vd", nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)

	assert.Equal(t, []string{"va", "vb", "vc"}, seenVals)
}

func TestWithQueryTransformerFunc(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	url := "/?a=va&b=vb&c=vc&d=vd"

	tests := []struct{
		name string
		keys []string
		seenVals []string
	}{
		{
			name: "some keys defined",
			keys: []string{"a", "b", "c"},
			seenVals: []string{"va", "vb", "vc"},
		},
		{
			name: "keys not in the query",
			keys: []string{"x", "y", "z"},
			seenVals: []string{},
		},
		{
			name: "no keys defined",
			keys: []string{},
			seenVals: []string{},
		},
		{
			name: "nil keys",
			keys: nil,
			seenVals: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seenVals := make([]string, 0)
			fn := func(s string) string {
				seenVals = append(seenVals, s)
				return s
			}

			h := middleware.WithQueryTransformerFunc(dummyHandler, tt.keys, fn)

			req, _ := http.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)

			assert.Equal(t, tt.seenVals, seenVals)
		})
	}

}
