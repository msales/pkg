package health

import (
	"github.com/go-zoo/bone"
	"github.com/msales/pkg/v3/httpx"
)

// DefaultPattern is the default health http path.
var DefaultPattern = "/health"

// DefaultHandler is the default health http handler.
var DefaultHandler = NewHandler()

// NewMux creates a Mux for the health endpoint
func NewMux(reporters ...Reporter) *bone.Mux {
	mux := httpx.NewMux()
	mux.Handle(DefaultPattern, DefaultHandler.With(reporters...))

	return mux
}
