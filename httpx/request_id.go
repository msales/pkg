package httpx

import "context"

// key used to store context values from within this package.
type key int

const (
	requestIDKey = iota
)

// WithRequestID inserts a RequestID into the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID extracts a RequestID from a context.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}
