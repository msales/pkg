package stats

import (
	"context"
	"time"
)

// DefaultHeartbeatInterval is the default heartbeat ticker interval.
var DefaultHeartbeatInterval = time.Second

// Heartbeat enters a loop, reporting a heartbeat counter periodically.
func Heartbeat(stats Stats) {
	HeartbeatEvery(stats, DefaultHeartbeatInterval)
}

// HeartbeatEvery enters a loop, reporting a heartbeat counter at the specified interval.
func HeartbeatEvery(stats Stats, t time.Duration) {
	c := time.Tick(t)
	for range c {
		_ = stats.Inc("heartbeat", 1, 1.0)
	}
}

// HeartbeatFromContext is the same as HeartbeatEvery but from context.
func HeartbeatFromContext(ctx context.Context, t time.Duration) {
	if s, ok := FromContext(ctx); ok {
		HeartbeatEvery(s, t)
	}
}
