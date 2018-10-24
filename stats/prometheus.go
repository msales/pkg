package stats

import (
	"errors"
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

// createKey creates a unique metric key.
func (s *Prometheus) createKey(name string, lblNames []string) string {
	return name + strings.Join(lblNames, ":")
}

// formatFQN formats FQN strings.
func (s *Prometheus) formatFQN(name string) string {
	return strings.Replace(name, ".", "_", -1)
}

// formatPrometheusTags create a prometheus Label map from tags.
func formatPrometheusTags(tags []interface{}) ([]string, prometheus.Labels) {
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