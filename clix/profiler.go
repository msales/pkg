package clix

import (
	"net/http"
)

func RunProfiler(c Ctx) {
	if !c.Bool(FlagProfiler) {
		return
	}

	listenAddress := ":" + c.String(FlagProfilerPort)

	go func() {

		err := http.ListenAndServe(listenAddress, nil)
		if err != nil {
			panic(err)
		}
	}()
}
