package clix

import (
	"time"

	"github.com/urfave/cli/v2"
)

// Flag constants declared for CLI use.
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

	FlagKafkaConsumerBrokers      = "kafka-consumer-brokers"
	FlagKafkaConsumerGroupID      = "kafka-consumer-group-id"
	FlagKafkaConsumerTopic        = "kafka-consumer-topic"
	FlagKafkaConsumerKafkaVersion = "kafka-consumer-kafka-version"
	FlagKafkaProducerBrokers      = "kafka-producer-brokers"
	FlagKafkaProducerTopic        = "kafka-producer-topic"
	FlagKafkaProducerKafkaVersion = "kafka-producer-kafka-version"

	FlagCommitBatch    = "commit-batch"
	FlagCommitInterval = "commit-interval"

	FlagRedisDSN = "redis-dsn"
)

type defaults struct {
	Port      string
	LogFormat string
	LogLevel  string

	ProfilerPort string
}

// Defaults holds the flag default values.
var Defaults = defaults{
	Port:      "80",
	LogFormat: "json",
	LogLevel:  "info",

	ProfilerPort: "8081",
}

// Flags represents a set of CLI flags.
type Flags []cli.Flag

// Merge joins one or more Flags together, making a new set.
func (f Flags) Merge(flags ...Flags) Flags {
	var m Flags
	m = append(m, f...)
	for _, flag := range flags {
		m = append(m, flag...)
	}

	return m
}

// ServerFlags are flags that configure a server.
var ServerFlags = Flags{
	&cli.StringFlag{
		Name:    FlagPort,
		Value:   Defaults.Port,
		Usage:   "Port for HTTP server to listen on",
		EnvVars: []string{"PORT"},
	},
}

// KafkaConsumerFlags are flags that configure a Kafka consumer.
var KafkaConsumerFlags = Flags{
	&cli.StringSliceFlag{
		Name:     FlagKafkaConsumerBrokers,
		Usage:    "Kafka consumer brokers.",
		EnvVars:  []string{"KAFKA_CONSUMER_BROKERS"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagKafkaConsumerGroupID,
		Usage:    "Kafka consumer group id.",
		EnvVars:  []string{"KAFKA_CONSUMER_GROUP_ID"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagKafkaConsumerTopic,
		Usage:    "Kafka topic to consume from.",
		EnvVars:  []string{"KAFKA_CONSUMER_TOPIC"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagKafkaConsumerKafkaVersion,
		Usage:    "Kafka version (e.g. 0.10.2.0 or 2.3.0).",
		EnvVars:  []string{"KAFKA_CONSUMER_KAFKA_VERSION"},
		Required: true,
	},
}

// KafkaProducerFlags are flags that configure a Kafka producer.
var KafkaProducerFlags = Flags{
	&cli.StringSliceFlag{
		Name:     FlagKafkaProducerBrokers,
		Usage:    "Kafka producer brokers.",
		EnvVars:  []string{"KAFKA_PRODUCER_BROKERS"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagKafkaProducerTopic,
		Usage:    "Kafka topic to produce into.",
		EnvVars:  []string{"KAFKA_PRODUCER_TOPIC"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     FlagKafkaProducerKafkaVersion,
		Usage:    "Kafka version (e.g. 0.10.2.0 or 2.3.0).",
		EnvVars:  []string{"KAFKA_PRODUCER_KAFKA_VERSION"},
		Required: true,
	},
}

// CommitterFlags are flags that configure message processing batch size and committing interval.
var CommitterFlags = Flags{
	&cli.IntFlag{
		Name:    FlagCommitBatch,
		Value:   500,
		Usage:   "Commit batch size for message processing.",
		EnvVars: []string{"COMMIT_BATCH"},
	},
	&cli.DurationFlag{
		Name:    FlagCommitInterval,
		Value:   1 * time.Second,
		Usage:   "Commit interval for message processing.",
		EnvVars: []string{"COMMIT_INTERVAL"},
	},
}

// RedisFlags are flags that configure redis.
var RedisFlags = Flags{
	&cli.StringFlag{
		Name:     FlagRedisDSN,
		Usage:    "The DSN of Redis.",
		EnvVars:  []string{"REDIS_DSN"},
		Required: true,
	},
}

// CommonFlags are flags that configure logging and stats.
var CommonFlags = Flags{
	&cli.StringFlag{
		Name:    FlagLogFormat,
		Value:   Defaults.LogFormat,
		Usage:   "Specify the format of logs. Supported formats: 'terminal', 'json'",
		EnvVars: []string{"LOG_FORMAT"},
	},
	&cli.StringFlag{
		Name:    FlagLogLevel,
		Value:   Defaults.LogLevel,
		Usage:   "Specify the log level. E.g. 'debug', 'warning'.",
		EnvVars: []string{"LOG_LEVEL"},
	},
	&cli.StringSliceFlag{
		Name:    FlagLogTags,
		Usage:   "A list of tags appended to every log. Format: key=value.",
		EnvVars: []string{"LOG_TAGS"},
	},
	&cli.StringFlag{
		Name:     FlagStatsDSN,
		Usage:    "The URL of a stats backend.",
		EnvVars:  []string{"STATS_DSN"},
		Required: true,
	},
	&cli.StringFlag{
		Name:    FlagStatsPrefix,
		Usage:   "The prefix of the measurements names.",
		EnvVars: []string{"STATS_PREFIX"},
	},
	&cli.StringSliceFlag{
		Name:    FlagStatsTags,
		Usage:   "A list of tags appended to every measurement. Format: key=value.",
		EnvVars: []string{"STATS_TAGS"},
	},
}

// ProfilerFlags are flags that configure to the profiler.
var ProfilerFlags = Flags{
	&cli.BoolFlag{
		Name:    FlagProfiler,
		Usage:   "Enable profiler server.",
		EnvVars: []string{"PROFILER"},
	},
	&cli.StringFlag{
		Name:    FlagProfilerPort,
		Value:   Defaults.ProfilerPort,
		Usage:   "Port for the profiler to listen on.",
		EnvVars: []string{"PROFILER_PORT"},
	},
}
