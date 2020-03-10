package stats

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus represents a promethus stats collector.
type Prometheus struct {
	prefix string

	reg      *prometheus.Registry
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
	timings  map[string]*prometheus.SummaryVec
}

// NewPrometheus creates a new Prometheus stats instance.
func NewPrometheus(prefix string) *Prometheus {
	return &Prometheus{
		prefix:   prefix,
		reg:      prometheus.NewRegistry(),
		counters: map[string]*prometheus.CounterVec{},
		gauges:   map[string]*prometheus.GaugeVec{},
		timings:  map[string]*prometheus.SummaryVec{},
	}
}

// Handler gets the prometheus HTTP handler.
func (s *Prometheus) Handler() http.Handler {
	return promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{})
}

// Inc increments a count by the value.
func (s *Prometheus) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	lblNames, lbls := formatPrometheusTags(tags)

	key := s.createKey(name, lblNames)
	m, ok := s.counters[key]
	if !ok {
		m = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: s.formatFQN(s.prefix),
				Name:      s.formatFQN(name),
				Help:      name,
			},
			lblNames,
		)

		err := s.reg.Register(m)
		if err == nil {
			s.counters[key] = m
		} else {
			existsErr, ok := err.(prometheus.AlreadyRegisteredError)
			if !ok {
				return err
			}

			m, ok = existsErr.ExistingCollector.(*prometheus.CounterVec)
			if !ok {
				return fmt.Errorf("stats: expected the collector to be instance of *CounterVec, got %T instead", existsErr.ExistingCollector)
			}
		}
	}

	m.With(lbls).Add(float64(value))

	return nil
}

// Dec decrements a count by the value.
func (s *Prometheus) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	return errors.New("prometheus: decrement not supported")
}

// Gauge measures the value of a metric.
func (s *Prometheus) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	lblNames, lbls := formatPrometheusTags(tags)

	key := s.createKey(name, lblNames)
	m, ok := s.gauges[key]
	if !ok {
		m = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: s.formatFQN(s.prefix),
				Name:      s.formatFQN(name),
				Help:      name,
			},
			lblNames,
		)

		err := s.reg.Register(m)
		if err == nil {
			s.gauges[key] = m
		} else {
			existsErr, ok := err.(prometheus.AlreadyRegisteredError)
			if !ok {
				return err
			}

			m, ok = existsErr.ExistingCollector.(*prometheus.GaugeVec)
			if !ok {
				return fmt.Errorf("stats: expected the collector to be instance of *GaugeVec, got %T instead", existsErr.ExistingCollector)
			}

		}
	}

	m.With(lbls).Set(value)

	return nil
}

// Timing sends the value of a Duration.
func (s *Prometheus) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	lblNames, lbls := formatPrometheusTags(tags)

	key := s.createKey(name, lblNames)
	m, ok := s.timings[key]
	if !ok {
		m = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  s.formatFQN(s.prefix),
				Name:       s.formatFQN(name),
				Help:       name,
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			lblNames,
		)

		err := s.reg.Register(m)
		if err == nil {
			s.timings[key] = m
		} else {
			existsErr, ok := err.(prometheus.AlreadyRegisteredError)
			if !ok {
				return err
			}

			m = existsErr.ExistingCollector.(*prometheus.SummaryVec)
			if !ok {
				return fmt.Errorf("stats: expected the collector to be instance of *SummaryVec, got %T instead", existsErr.ExistingCollector)
			}

		}
	}

	m.With(lbls).Observe(float64(value) / float64(time.Millisecond))

	return nil
}

// Close closes the client and flushes buffered stats, if applicable
func (s *Prometheus) Close() error {
	return nil
}

// createKey creates a unique metric key.
func (s *Prometheus) createKey(name string, lblNames []string) string {
	return name + strings.Join(lblNames, ":")
}

// formatFQN formats FQN strings.
func (s *Prometheus) formatFQN(name string) string {
	r := strings.NewReplacer(".", "_", "-", "_")

	return r.Replace(name)
}

// formatPrometheusTags create a prometheus Label map from tags.
func formatPrometheusTags(tags []interface{}) ([]string, prometheus.Labels) {
	tags = deduplicateTags(normalizeTags(tags))

	b := make([]byte, 0, 65) // The largest needed buffer is 65 bytes for a signed int64.

	names := make([]string, 0, len(tags)/2)
	lbls := make(prometheus.Labels, len(tags)/2)
	for i := 0; i < len(tags); i += 2 {
		key, ok := tags[i].(string) // The stats key must be a string.
		if !ok {
			key = fmt.Sprintf("STATS_ERROR: key %v is not a string", tags[i])
		}
		names = append(names, key)

		lbl, ok := toString(tags[i+1], b)
		if !ok {
			lbl = string(b)
		}
		lbls[key] = lbl
	}

	return names, lbls
}

// toString converts the given value to a string. It either returns the new string and true
// or fills the passed byte slice and returns an empty string and false. The user needs to check
// the returned boolean and take the string (if true) or get data from the slice.
// This is the optimization: filling the buffer allows to re-use the memory and avoid
// allocations when converting floats. Returning the string directly avoids copying strings.
func toString(v interface{}, b []byte) (string, bool) {
	switch vv := v.(type) {
	case string:
		return vv, true
	case bool:
		strconv.AppendBool(b, vv)
	case float32:
		strconv.AppendFloat(b, float64(vv), 'f', -1, 64)
	case float64:
		strconv.AppendFloat(b, vv, 'f', -1, 64)
	case int:
		strconv.AppendInt(b, int64(vv), 10)
	case int8:
		strconv.AppendInt(b, int64(vv), 10)
	case int16:
		strconv.AppendInt(b, int64(vv), 10)
	case int32:
		strconv.AppendInt(b, int64(vv), 10)
	case int64:
		strconv.AppendInt(b, vv, 10)
	case uint:
		strconv.AppendUint(b, uint64(vv), 10)
	case uint8:
		strconv.AppendUint(b, uint64(vv), 10)
	case uint16:
		strconv.AppendUint(b, uint64(vv), 10)
	case uint32:
		strconv.AppendUint(b, uint64(vv), 10)
	case uint64:
		strconv.AppendUint(b, vv, 10)
	default:
		return fmt.Sprintf("STATS_ERROR: cannot convert value %v to string", vv), true
	}

	return "", false
}
