package middleware

import (
	"context"
	"net/http"
)

// WithContext set the context on the request.
func WithContext(h http.Handler, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
