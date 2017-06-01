package stats

import "testing"

func TestRuntimeStats(t *testing.T) {
	runtime := newRuntimeStats()
	runtime.send(Null)
}
