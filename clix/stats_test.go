package clix_test

import (
	"testing"

	"github.com/msales/pkg/v4/clix"
	"github.com/msales/pkg/v4/log"
	"github.com/msales/pkg/v4/stats"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		dsn    string
		prefix string
		tags   *cli.StringSlice

		shouldErr bool
	}{
		{"statsd://localhost:8125", "test", cli.NewStringSlice(), false},
		{"", "test", cli.NewStringSlice(), false},
		{"l2met://", "test", cli.NewStringSlice(), false},
		{"prometheus://", "test", cli.NewStringSlice(), false},
		{"prometheus://:51234", "test", cli.NewStringSlice(), false},
		{"l2met://", "", cli.NewStringSlice(), false},
		{"invalid-scheme", "", cli.NewStringSlice(), true},
		{"unknownscheme://", "", cli.NewStringSlice(), true},
		{":/", "", cli.NewStringSlice(), true},
		{"l2met://", "", cli.NewStringSlice("a"), true},
	}

	for _, tt := range tests {
		c, fs := newTestContext()
		fs.String(clix.FlagStatsDSN, tt.dsn, "doc")
		fs.String(clix.FlagStatsPrefix, tt.prefix, "doc")
		fs.Var(tt.tags, clix.FlagStatsTags, "doc")

		s, err := clix.NewStats(c, log.Null)

		if tt.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Implements(t, (*stats.Stats)(nil), s)
		}
	}
}
