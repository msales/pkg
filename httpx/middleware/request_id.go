package middleware

import (
	"net/http"

	"github.com/msales/pkg/httpx"
)

// DefaultRequestIDExtractor is the default function to use to extract a request
// id from an http.Request.
var DefaultRequestIdExtractor = httpx.HeaderExtractor([]string{"X-Request-Id", "Request-Id"})

// RequestID is a middleware that extracts the request id from the request
// and inserts it into context.
type RequestID struct {
	// Extractor is a function that can extract a request id from an http.Request.
	Extractor func(*http.Request) string

	handler http.Handler
}

// ExtractRequestID extracts the request id from the request.
func ExtractRequestID(h http.Handler) *RequestID {
	return &RequestID{
		handler: h,
	}
}

// ServeHTTP serves the request.
func (h RequestID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := h.Extractor
	if e == nil {
		e = DefaultRequestIdExtractor
	}
	requestID := e(r)

	ctx := httpx.WithRequestID(r.Context(), requestID)
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}
