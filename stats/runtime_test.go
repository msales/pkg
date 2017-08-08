package stats_test

import (
	"testing"
	"time"

	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestRuntime(t *testing.T) {
	m := new(MockStats)
	m.On("Gauge", "runtime.cpu.goroutines", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	stats.DefaultRuntimeInterval = 10 * time.Microsecond

	go stats.Runtime(m)

	time.Sleep(time.Millisecond)

	m.AssertExpectations(t)
}
