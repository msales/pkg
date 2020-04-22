package clix_test

import (
	"testing"

	"github.com/msales/pkg/v4/clix"
	"github.com/msales/pkg/v4/log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level  string
		format string
		tags   *cli.StringSlice

		shouldErr bool
	}{
		{"info", "json", cli.NewStringSlice(), false},
		{"info", "terminal", cli.NewStringSlice(), false},
		{"info", "logfmt", cli.NewStringSlice(), false},
		{"", "json", cli.NewStringSlice(), false},
		{"info", "", cli.NewStringSlice(), false},
		{"invalid", "json", cli.NewStringSlice(), true},
		{"info", "invalid", cli.NewStringSlice(), true},
		{"info", "json", cli.NewStringSlice("string"), true},
	}

	for _, tt := range tests {
		c, fs := newTestContext()
		fs.String(clix.FlagLogLevel, tt.level, "doc")
		fs.String(clix.FlagLogFormat, tt.format, "doc")
		fs.Var(tt.tags, clix.FlagLogTags, "doc")

		l, err := clix.NewLogger(c)

		if tt.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Implements(t, (*log.Logger)(nil), l)
		}
	}
}
