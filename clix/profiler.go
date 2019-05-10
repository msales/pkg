package clix

import (
	"context"
	"net/http"
	"net/http/pprof"

	"gopkg.in/urfave/cli.v1"
)

var profilerServer = &http.Server{}

// RunProfiler runs a profiler server.
func RunProfiler(c *cli.Context) error {
	if !c.GlobalBool(FlagProfiler) {
		return nil
	}

	profilerServer.Handler = newProfilerMux()
	profilerServer.Addr = ":" + c.GlobalString(FlagProfilerPort)

	go func() {
		err := profilerServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

// StopProfiler stops a running profiler server.
func StopProfiler() error {
	return profilerServer.Shutdown(context.Background())
}

func newProfilerMux() http.Handler {
	mux := &http.ServeMux{}

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}
