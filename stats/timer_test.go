package stats_test

import (
	"context"
	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTimer(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", mock.Anything, float32(1.0), mock.Anything).Return(nil)

	ctx := stats.WithStats(context.Background(), m)
	ti := stats.Time(ctx, "test", 1.0, nil)
	ti.Done()

	m.AssertExpectations(t)
}
