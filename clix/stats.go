package clix

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"gopkg.in/urfave/cli.v1"
)

// NewStats creates a new stats client.
func NewStats(c *cli.Context, l log.Logger) (stats.Stats, error) {
	var s stats.Stats
	var err error

	dsn := c.GlobalString(FlagStatsDSN)
	if dsn == "" {
		return stats.Null, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch scheme := uri.Scheme; scheme {
	case "statsd":
		s, err = newStatsDStats(c, uri.Host)
		if err != nil {
			return nil, err
		}

	case "l2met":
		s = newL2metStats(c, l)

	case "prometheus":
		s = newPrometheusStats(c, uri.Host, l)

	default:
		return nil, fmt.Errorf("Unknown scheme: %s", scheme)
	}

	tags, err := SplitTags(c.GlobalStringSlice(FlagStatsTags), "=")
	if err != nil {
		return nil, err
	}

	return stats.NewTaggedStats(s, tags...), nil
}

func newStatsDStats(c *cli.Context, addr string) (stats.Stats, error) {
	s, err := stats.NewBufferedStatsd(addr, c.GlobalString(FlagStatsPrefix), stats.WithFlushInterval(1*time.Second))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func newL2metStats(c *cli.Context, l log.Logger) stats.Stats {
	return stats.NewL2met(l, c.GlobalString(FlagStatsPrefix))
}

func newPrometheusStats(c *cli.Context, addr string, l log.Logger) stats.Stats {
	s := stats.NewPrometheus(c.GlobalString(FlagStatsPrefix))

	if addr != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", s.Handler())
		go func() {
			if err := http.ListenAndServe(addr, mux); err != nil {
				l.Error(err.Error())
			}
		}()
	}

	return s
}
