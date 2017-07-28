package stats_test

import (
	"testing"
	"github.com/msales/pkg/stats"
	"context"
	"github.com/stretchr/testify/mock"
)

func TestTimer(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", mock.Anything, float32(1.0), mock.Anything).Return(nil)

	ctx := stats.WithStats(context.Background(), m)
	ti := stats.Time(ctx, "test", 1.0, nil)
	ti.Done()

	m.AssertExpectations(t)
}
