package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/msales/pkg/v3/stats"
	"github.com/stretchr/testify/mock"
)

func TestHeartbeat(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "heartbeat", int64(1), float32(1.0), mock.Anything).Return(nil)
	stats.DefaultHeartbeatInterval = time.Millisecond

	go stats.Heartbeat(m)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Inc", "heartbeat", int64(1), float32(1.0), mock.Anything)
}

func TestHeartbeatFromContext(t *testing.T) {
	m := new(MockStats)
	m.On("Inc", "heartbeat", int64(1), float32(1.0), mock.Anything).Return(nil)
	ctx := stats.WithStats(context.Background(), m)

	go stats.HeartbeatFromContext(ctx, time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	m.AssertCalled(t, "Inc", "heartbeat", int64(1), float32(1.0), mock.Anything)
}
