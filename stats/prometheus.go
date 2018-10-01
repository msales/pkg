package stats

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus represents a promethus stats collector.
type Prometheus struct {
	prefix string

	reg      *prometheus.Registry
	metrics  map[string]prometheus.Collector
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
	timings  map[string]*prometheus.SummaryVec
}

// NewPrometheus creates a new Prometheus stats instance.
func NewPrometheus(prefix string) Stats {
	return &Prometheus{
		prefix: prefix,
		reg:    prometheus.NewRegistry(),
	}
}

// Handler gets the prometheus HTTP handler.
func (s *Prometheus) Handler() http.Handler {
	return promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{})
}

// Inc increments a count by the value.
func (s *Prometheus) Inc(name string, value int64, rate float32, tags ...interface{}) error {
	lblNames, lbls := prometheusTags(tags)
	key := prometheusKey(name, lblNames)

	m, ok := s.counters[key]
	if !ok {
		m := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: s.prefix,
				Name:      name,
			},
			lblNames,
		)

		if err := s.reg.Register(m); err != nil {
			return err
		}
		s.counters[key] = m
	}

	m.With(lbls).Add(float64(value))

	return nil
}

// Dec decrements a count by the value.
func (s *Prometheus) Dec(name string, value int64, rate float32, tags ...interface{}) error {
	lblNames, lbls := prometheusTags(tags)
	key := prometheusKey(name, lblNames)

	m, ok := s.counters[key]
	if !ok {
		m := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: s.prefix,
				Name:      name,
			},
			lblNames,
		)

		if err := s.reg.Register(m); err != nil {
			return err
		}
		s.counters[key] = m
	}

	m.With(lbls).Add(-float64(value))

	return nil
}

// Gauge measures the value of a metric.
func (s *Prometheus) Gauge(name string, value float64, rate float32, tags ...interface{}) error {
	lblNames, lbls := prometheusTags(tags)
	key := prometheusKey(name, lblNames)

	m, ok := s.gauges[key]
	if !ok {
		m := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: s.prefix,
				Name:      name,
			},
			lblNames,
		)

		if err := s.reg.Register(m); err != nil {
			return err
		}
		s.gauges[key] = m
	}

	m.With(lbls).Set(value)

	return nil
}

// Timing sends the value of a Duration.
func (s *Prometheus) Timing(name string, value time.Duration, rate float32, tags ...interface{}) error {
	lblNames, lbls := prometheusTags(tags)
	key := prometheusKey(name, lblNames)

	m, ok := s.timings[key]
	if !ok {
		m := prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  s.prefix,
				Name:       name,
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			lblNames,
		)

		if err := s.reg.Register(m); err != nil {
			return err
		}
		s.timings[key] = m
	}

	m.With(lbls).Observe(float64(value) / float64(time.Millisecond))

	return nil
}

// Close closes the client and flushes buffered stats, if applicable
func (s *Prometheus) Close() error {
	return nil
}

// prometheusTags create a prometheus Label map from tags
func prometheusTags(tags []interface{}) ([]string, prometheus.Labels) {
	tags = deduplicateTags(normalizeTags(tags))

	names := make([]string, 0, len(tags)/2)
	lbls := make(prometheus.Labels, len(tags)/2)
	for i := 0; i < len(tags); i += 2 {
		key := fmt.Sprintf("%v", tags[i])
		names = append(names, key)
		lbls[key] = fmt.Sprintf("%v", tags[i+1])
	}

	sort.Strings(names)

	return names, lbls
}

// prometheusKey creates a unique metric key.
func prometheusKey(name string, lblNames []string) string {
	return name + strings.Join(lblNames, ":")
}
