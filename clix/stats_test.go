package clix_test

import (
	"testing"

	"github.com/msales/pkg/v2/clix"
	"github.com/msales/pkg/v2/log"
	"github.com/msales/pkg/v2/stats"
	"github.com/stretchr/testify/assert"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		dsn    string
		prefix string
		tags   []string

		shouldErr bool
	}{
		{"statsd://localhost:8125", "test", []string{}, false},
		{"", "test", []string{}, false},
		{"l2met://", "test", []string{}, false},
		{"l2met://", "", []string{}, false},
		{"invalid-scheme", "", []string{}, true},
		{"unknownscheme://", "", []string{}, true},
		{"l2met://", "", []string{"a"}, true},
	}

	for _, tt := range tests {
		ctx := new(CtxMock)
		ctx.On("String", clix.FlagStatsDSN).Return(tt.dsn)
		ctx.On("String", clix.FlagStatsPrefix).Return(tt.prefix)
		ctx.On("StringSlice", clix.FlagStatsTags).Return(tt.tags)

		s, err := clix.NewStats(ctx, log.Null)

		if tt.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Implements(t, (*stats.Stats)(nil), s)
		}
	}
}
