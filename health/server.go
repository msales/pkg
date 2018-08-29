package health

import (
	"net/http"
	"time"
)

const healthPath = "/health"

func NewServer(addr string, reporters ...Reporter) *http.Server {
	mux := &http.ServeMux{}
	mux.Handle(healthPath, DefaultHandler.With(reporters...))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
