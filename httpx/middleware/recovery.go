package middleware

import (
	"fmt"
	"net/http"

	"github.com/msales/pkg/v3/log"
)

// Recovery is a middleware that will recover from panics and logs the error.
type Recovery struct {
	handler http.Handler
}

// WithRecovery recovers from panics and log the error.
func WithRecovery(h http.Handler) http.Handler {
	return &Recovery{
		handler: h,
	}
}

// ServeHTTP serves the request.
func (m Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("%v", v)
			if v, ok := v.(error); ok {
				err = v
			}

			log.Error(r.Context(), err.Error())
			w.WriteHeader(500)
		}
	}()

	m.handler.ServeHTTP(w, r)
}
