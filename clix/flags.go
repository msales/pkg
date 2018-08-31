package clix

import (
	"errors"
	"strings"

	"github.com/urfave/cli"
)

const (
	FlagPort = "port"

	FlagLogFormat = "log-format"
	FlagLogLevel  = "log-level"
	FlagLogTags   = "log-tags"

	FlagStatsDSN    = "stats-dsn"
	FlagStatsPrefix = "stats-prefix"
	FlagStatsTags   = "stats-tags"

	FlagProfiler     = "profiler"
	FlagProfilerPort = "profiler-port"
)

type defaults struct {
	Port      int
	LogFormat string
	LogLevel  string

	ProfilerPort int
}

var Defaults = defaults{
	Port:      80,
	LogFormat: "json",
	LogLevel:  "info",

	ProfilerPort: 8081,
}

type Flags []cli.Flag

func (f Flags) Merge(flags Flags) Flags {
	return append(f, flags...)
}

var ServerFlags = Flags{
	cli.IntFlag{
		Name:   FlagPort,
		Value:  Defaults.Port,
		Usage:  "Port for HTTP server to listen on",
		EnvVar: "PORT",
	},
}

var CommonFlags = Flags{
	cli.StringFlag{
		Name:   FlagLogFormat,
		Value:  Defaults.LogFormat,
		Usage:  "Specify the format of logs. Supported formats: 'terminal', 'json'",
		EnvVar: "LOG_FORMAT",
	},
	cli.StringFlag{
		Name:   FlagLogLevel,
		Value:  Defaults.LogLevel,
		Usage:  "Specify the log level. E.g. `debug`, `warning`.",
		EnvVar: "LOG_LEVEL",
	},
	cli.StringSliceFlag{
		Name:   FlagLogTags,
		Usage:  "A list of tags appended to every log. Format: key=value.",
		EnvVar: "LOG_TAGS",
	},

	cli.StringFlag{
		Name:   FlagStatsDSN,
		Usage:  "The URL of a stats backend.",
		EnvVar: "STATS_DSN",
	},
	cli.StringFlag{
		Name:   FlagStatsPrefix,
		Usage:  "The prefix of the measurements names.",
		EnvVar: "STATS_PREFIX",
	},
	cli.StringSliceFlag{
		Name:   FlagStatsTags,
		Usage:  "A list of tags appended to every measurement. Format: key=value.",
		EnvVar: "STATS_TAGS",
	},
}

var ProfilerFlags = Flags{
	cli.BoolFlag{
		Name:   FlagProfiler,
		Usage:  "Enable profiler server.",
		EnvVar: "PROFILER",
	},
	cli.IntFlag{
		Name:   FlagProfilerPort,
		Value:  Defaults.ProfilerPort,
		Usage:  "Port for the profiler to listen on.",
		EnvVar: "PROFILER_PORT",
	},
}

func SplitTags(slice []string, sep string) ([]interface{}, error) {
	res := make([]interface{}, 2*len(slice))

	for i, str := range slice {
		parts := strings.SplitN(str, sep, 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid tags string")
		}

		res[2*i] = parts[0]
		res[2*i+1] = parts[1]
	}

	return res, nil
}
