package clix_test

import (
	"testing"

	"github.com/msales/pkg/clix"
	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		dsn    string
		prefix string
		tags   *cli.StringSlice

		shouldErr bool
	}{
		{"statsd://localhost:8125", "test", &cli.StringSlice{}, false},
		{"", "test", &cli.StringSlice{}, false},
		{"l2met://", "test", &cli.StringSlice{}, false},
		{"prometheus://", "test", &cli.StringSlice{}, false},
		{"prometheus://:51234", "test", &cli.StringSlice{}, false},
		{"l2met://", "", &cli.StringSlice{}, false},
		{"invalid-scheme", "", &cli.StringSlice{}, true},
		{"unknownscheme://", "", &cli.StringSlice{}, true},
		{"l2met://", "", &cli.StringSlice{"a"}, true},
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
