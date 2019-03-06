package middleware

import (
	"context"
	"net/http"
)

// WithCommon wraps the handler with the commonly used middlewares.
func WithCommon(ctx context.Context, h http.Handler, fns ...TagsFunc) http.Handler {
	h = WithResponseTime(h, fns...) // Innermost
	h = WithRequestStats(h, fns...)
	h = WithRecovery(h)
	h = WithContext(ctx, h) // Outermost

	return h
}
