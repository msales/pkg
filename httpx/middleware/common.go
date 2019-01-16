package middleware

import (
	"context"
	"net/http"
)

// WithCommon wraps the handler with the commonly used middlewares.
func WithCommon(ctx context.Context, h http.Handler) http.Handler {
	h = WithResponseTime(h) // Innermost
	h = WithRequestStats(h)
	h = WithContext(ctx, h)
	h = WithRecovery(h) // Outermost

	return h
}
