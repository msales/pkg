package clix

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/go-zoo/bone"
	"github.com/msales/pkg/v3/httpx"
	"gopkg.in/urfave/cli.v1"
)

var profilerServer = &http.Server{}

// RunProfiler runs a profiler server.
func RunProfiler(c *cli.Context) error {
	if !c.Bool(FlagProfiler) {
		return nil
	}

	profilerServer.Handler = newProfilerMux()
	profilerServer.Addr = ":" + c.String(FlagProfilerPort)

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

func newProfilerMux() *bone.Mux {
	mux := httpx.NewMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}
