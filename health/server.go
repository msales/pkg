package health

import (
	"net/http"
	"time"
)

var DefaultPattern = "/health"

func ListenAndServe(addr string, reporters ...Reporter) error {
	return newServer(addr, reporters...).ListenAndServe()
}

func newServer(addr string, reporters ...Reporter) *http.Server {
	mux := &http.ServeMux{}
	mux.Handle(DefaultPattern, DefaultHandler.With(reporters...))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
