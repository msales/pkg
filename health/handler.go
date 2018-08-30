package health

import (
	"net/http"
)

var DefaultHandler = &Handler{}

type Reporter interface {
	IsHealthy() error
}

type Handler struct {
	Reporters []Reporter
	ShowErr   bool
}

func (h *Handler) With(reporters ...Reporter) *Handler {
	h.Reporters = append(h.Reporters, reporters...)
	return h
}

func (h *Handler) WithErrors() *Handler {
	h.ShowErr = true
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, reporter := range h.Reporters {
		if err := reporter.IsHealthy(); err != nil {
			http.Error(w, h.getResponseContent(err), http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getResponseContent(err error) string {
	if h.ShowErr {
		return err.Error()
	}

	return http.StatusText(http.StatusServiceUnavailable)
}
