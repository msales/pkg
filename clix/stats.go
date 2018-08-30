package clix

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
)

func NewStats(c Ctx, l log.Logger) (stats.Stats, error) {
	var s stats.Stats
	var err error

	dsn := c.String(FlagStatsAddr)
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
	default:
		return nil, errors.New(fmt.Sprintf("Unknown scheme: %s", scheme))
	}

	tags, err := splitTags(c.StringSlice(FlagStatsTags), "=")
	if err != nil {
		return nil, err
	}

	return stats.NewTaggedStats(s, tags...), nil
}

func newStatsDStats(c Ctx, addr string) (stats.Stats, error) {
	s, err := stats.NewBufferedStatsd(addr, c.String(FlagStatsPrefix), stats.WithFlushInterval(1*time.Second))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func newL2metStats(c Ctx, l log.Logger) stats.Stats {
	return stats.NewL2met(l, c.String(FlagStatsPrefix))
}
