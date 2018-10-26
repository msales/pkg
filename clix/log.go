package clix

import (
	"fmt"
	"os"

	"github.com/inconshreveable/log15"
)

// NewLogger creates a new logger.
func NewLogger(c *cli.Context) (log15.Logger, error) {
	levelStr := c.String(FlagLogLevel)
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := log15.LvlFromString(levelStr)
	if err != nil {
		return nil, err
	}

	format, err := newLogFormat(c)
	if err != nil {
		return nil, err
	}

	handler := log15.LvlFilterHandler(level, log15.StreamHandler(os.Stdout, format))

	tags, err := SplitTags(c.StringSlice(FlagLogTags), "=")
	if err != nil {
		return nil, err
	}

	logger := log15.New(tags...)
	logger.SetHandler(handler)

	return logger, nil
}

func newLogFormat(c *cli.Context) (log15.Format, error) {
	format := c.String(FlagLogFormat)
	switch format {
	case "terminal":
		return log15.TerminalFormat(), nil
	case "json", "":
		return log15.JsonFormat(), nil
	case "logfmt":
		return log15.LogfmtFormat(), nil
	default:
		return nil, fmt.Errorf("invalid log format: '%s'", format)
	}
}
