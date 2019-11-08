package clix

import (
	"time"

	"gopkg.in/urfave/cli.v1"
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
	cli.StringFlag{
		Name:   FlagPort,
		Value:  Defaults.Port,
		Usage:  "Port for HTTP server to listen on",
		EnvVar: "PORT",
	},
}

// KafkaConsumerFlags are flags that configure a Kafka consumer.
var KafkaConsumerFlags = Flags{
	cli.StringSliceFlag{
		Name:   FlagKafkaConsumerBrokers,
		Usage:  "Kafka consumer brokers.",
		EnvVar: "KAFKA_CONSUMER_BROKERS",
	},
	cli.StringFlag{
		Name:   FlagKafkaConsumerGroupID,
		Usage:  "Kafka consumer group id.",
		EnvVar: "KAFKA_CONSUMER_GROUP_ID",
	},
	cli.StringFlag{
		Name:   FlagKafkaConsumerTopic,
		Usage:  "Kafka topic to consume from.",
		EnvVar: "KAFKA_CONSUMER_TOPIC",
	},
	cli.StringFlag{
		Name: FlagKafkaConsumerKafkaVersion,
		Usage: "Kafka version (e.g. 0.10.2.0 or 2.3.0).",
		EnvVar: "KAFKA_CONSUMER_KAFKA_VERSION",
	},
}

// KafkaProducerFlags are flags that configure a Kafka producer.
var KafkaProducerFlags = Flags{
	cli.StringSliceFlag{
		Name:   FlagKafkaProducerBrokers,
		Usage:  "Kafka producer brokers.",
		EnvVar: "KAFKA_PRODUCER_BROKERS",
	},
	cli.StringFlag{
		Name:   FlagKafkaProducerTopic,
		Usage:  "Kafka topic to produce into.",
		EnvVar: "KAFKA_PRODUCER_TOPIC",
	},
	cli.StringFlag{
		Name: FlagKafkaProducerKafkaVersion,
		Usage: "Kafka version (e.g. 0.10.2.0 or 2.3.0).",
		EnvVar: "KAFKA_PRODUCER_KAFKA_VERSION",
	},
}

// CommitterFlags are flags that configure message processing batch size and committing interval.
var CommitterFlags = Flags{
	cli.IntFlag{
		Name:   FlagCommitBatch,
		Value:  500,
		Usage:  "Commit batch size for message processing.",
		EnvVar: "COMMIT_BATCH",
	},
	cli.DurationFlag{
		Name:   FlagCommitInterval,
		Value:  1 * time.Second,
		Usage:  "Commit interval for message processing.",
		EnvVar: "COMMIT_INTERVAL",
	},
}

// RedisFlags are flags that configure redis.
var RedisFlags = Flags{
	cli.StringFlag{
		Name:   FlagRedisDSN,
		Usage:  "The DSN of Redis.",
		EnvVar: "REDIS_DSN",
	},
}

// CommonFlags are flags that configure logging and stats.
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
		Usage:  "Specify the log level. E.g. 'debug', 'warning'.",
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

// ProfilerFlags are flags that configure to the profiler.
var ProfilerFlags = Flags{
	cli.BoolFlag{
		Name:   FlagProfiler,
		Usage:  "Enable profiler server.",
		EnvVar: "PROFILER",
	},
	cli.StringFlag{
		Name:   FlagProfilerPort,
		Value:  Defaults.ProfilerPort,
		Usage:  "Port for the profiler to listen on.",
		EnvVar: "PROFILER_PORT",
	},
}
