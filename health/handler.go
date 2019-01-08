package health

import (
	"net/http"
)

// Reporter represents an a health reporter.
type Reporter interface {
	// IsHealthy emits error if application is not healthy.
	IsHealthy() error
}

// ReporterFunc is an adapter for anonymous functions to be used as health reporters.
type ReporterFunc func() error

// IsHealthy emits error if application is not healthy.
func (f ReporterFunc) IsHealthy() error {
	return f()
}

// Handler is an http health handler.
type Handler struct {
	reporters []Reporter
	showErr   bool
}

// NewHandler creates a new Handler instance.
func NewHandler() *Handler {
	return &Handler{}
}

// With adds reports to the handler.
func (h *Handler) With(reporters ...Reporter) *Handler {
	h.reporters = append(h.reporters, reporters...)
	return h
}

// WithErrors configures the handler to show the error message
// in the response.
func (h *Handler) WithErrors() *Handler {
	h.showErr = true
	return h
}

// ServeHTTP serves an http request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, reporter := range h.reporters {
		if err := reporter.IsHealthy(); err != nil {
			http.Error(w, h.getResponseContent(err), http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getResponseContent(err error) string {
	if h.showErr {
		return err.Error()
	}

	return http.StatusText(http.StatusServiceUnavailable)
}
