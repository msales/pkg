package middleware

import (
	"context"
	"net/http"
)

// WithCommon wraps the handler with the commonly used middlewares.
func WithCommon(ctx context.Context, h http.Handler) http.Handler {
	h = WithResponseTime(h) // Innermost
	h = WithRequestStats(h)
	h = WithRecovery(h)
	h = WithContext(ctx, h) // Outermost

	return h
}
