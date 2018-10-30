package health

import (
	"context"
	"net/http"
	"time"
)

// DefaultPattern is the default health http path.
var DefaultPattern = "/health"

// DefaultHandler is the default health http handler.
var DefaultHandler = NewHandler()

var server = &http.Server{
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
}

// StartServer starts the http health server on the given port
// with the given reporters.
func StartServer(addr string, reporters ...Reporter) error {
	server.Addr = addr
	server.Handler = newMux(reporters)

	return server.ListenAndServe()
}

// StopServer stops the htt[ health server
func StopServer() error {
	return server.Shutdown(context.Background())
}

func newMux(reporters []Reporter) http.Handler {
	mux := &http.ServeMux{}
	mux.Handle(DefaultPattern, DefaultHandler.With(reporters...))

	return mux
}
