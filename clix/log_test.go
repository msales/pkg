package clix_test

import (
	"testing"

	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level  string
		format string
		tags   *cli.StringSlice

		shouldErr bool
	}{
		{"info", "json", &cli.StringSlice{}, false},
		{"info", "terminal", &cli.StringSlice{}, false},
		{"info", "logfmt", &cli.StringSlice{}, false},
		{"", "json", &cli.StringSlice{}, false},
		{"info", "", &cli.StringSlice{}, false},
		{"invalid", "json", &cli.StringSlice{}, true},
		{"info", "invalid", &cli.StringSlice{}, true},
		{"info", "json", &cli.StringSlice{"single"}, true},
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
