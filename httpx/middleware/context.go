package middleware

import (
	"context"
	"net/http"
)

// WithContext set the context on the request.
func WithContext(ctx context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(withValueContext(r.Context(), ctx)))
	})
}

// mergedContext represents a context that merges a parent with a value context.
type mergedContext struct {
	context.Context

	valueCtx context.Context
}

// mergeContexts returns a new context which is child of the two passed contexts.
//
// Done() returns the channel from the parent context.
//
// Deadline() returns the deadline from the parent context.
//
// Err() returns the error from the parent context.
//
// Value(key) looks for key in the value context first and falls back to the parent.
func withValueContext(parent, valueCtx context.Context) context.Context {
	return &mergedContext{parent, valueCtx}
}

// Value returns the value associated with this context for key, or nil.
func (c *mergedContext) Value(key interface{}) interface{} {
	v := c.valueCtx.Value(key)
	if v != nil {
		return v
	}
	return c.Context.Value(key)
}
