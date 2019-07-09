package middleware

import (
	"fmt"
	"github.com/msales/pkg/v3/stats"
	"net/http"
	"runtime/debug"

	"github.com/msales/pkg/v3/log"
)


// RecoveryFunc is used to configure the recovery handler.
type RecoveryFunc func(*Recovery)

// WithoutStack disables the stack trace dump from the recovery log.
func WithoutStack() RecoveryFunc {
	return func(r *Recovery) {
		r.withStack = false
	}
}

// Recovery is a middleware that will recover from panics and logs the error.
type Recovery struct {
	handler   http.Handler
	withStack bool
}

// WithRecovery recovers from panics and log the error.
func WithRecovery(h http.Handler, opts ...RecoveryFunc) http.Handler {
	r := &Recovery{
		handler:   h,
		withStack: true,
	}

	for _, fn := range opts {
		fn(r)
	}

	return r
}

// ServeHTTP serves the request.
func (m *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("%v", v)
			if v, ok := v.(error); ok {
				err = v
			}

			var logCtx []interface{}
			if m.withStack {
				logCtx = append(logCtx, "stack", string(debug.Stack()))
			}

			log.Error(r.Context(), err.Error(), logCtx...)
			_ = stats.Inc(r.Context(), "panic_recovery", 1, 1)
			w.WriteHeader(500)
		}
	}()

	m.handler.ServeHTTP(w, r)
}
