package clix

import (
	"fmt"
	"os"
	"time"

	"github.com/msales/logged"
	"gopkg.in/urfave/cli.v1"
)

// NewLogger creates a new logger.
func NewLogger(c *cli.Context) (logged.Logger, error) {
	levelStr := c.String(FlagLogLevel)
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := logged.LevelFromString(levelStr)
	if err != nil {
		return nil, err
	}

	format, err := newLogFormat(c)
	if err != nil {
		return nil, err
	}

	h := logged.BufferedStreamHandler(os.Stdout, 2000, 1*time.Second, format)
	h = logged.LevelFilterHandler(level, h)

	tags, err := SplitTags(c.StringSlice(FlagLogTags), "=")
	if err != nil {
		return nil, err
	}

	logger := logged.New(h, tags...)

	return logger, nil
}

func newLogFormat(c *cli.Context) (logged.Formatter, error) {
	format := c.String(FlagLogFormat)
	switch format {
	case "terminal":
		fmt.Println("clix: terminal format depreciated, using logfmt instead")
		return logged.LogfmtFormat(), nil
	case "json", "":
		return logged.JSONFormat(), nil
	case "logfmt":
		return logged.LogfmtFormat(), nil
	default:
		return nil, fmt.Errorf("invalid log format: '%s'", format)
	}
}
