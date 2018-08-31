package clix

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"
)

var ProfilerServer = &http.Server{
	ReadTimeout: time.Minute,
	WriteTimeout: time.Minute,
}

func RunProfiler(c Ctx) {
	if !c.Bool(FlagProfiler) {
		return
	}

	ProfilerServer.Handler = makeProfilerMux()
	ProfilerServer.Addr = fmt.Sprintf(":%d", c.Int(FlagProfilerPort))

	go func() {
		err := ProfilerServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func makeProfilerMux() http.Handler {
	mux := &http.ServeMux{}

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}