package health

import (
	"fmt"
	"net/http"
)

type Checker interface {
	CheckHealth() error
}

type Handler struct {
	Checkers []Checker
}

func NewHandler(checkers ...Checker) *Handler {
	return &Handler{Checkers: checkers}
}

func (h *Handler) With(checker Checker) {
	h.Checkers = append(h.Checkers, checker)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, checker := range h.Checkers {
		if err := checker.CheckHealth(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
