package umami_prometheus

// Integration with Prometheus golang client **V2**
// No support for V1 exists for now

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/SimonDaKappa/go-umami"
)

const (
	PrometheusBackendName string = "prometheus"
)

// Mock backend for demonstration
type prometheusBackend struct {
	registry *prometheus.Registry
}

func NewPrometheusBackend(reg *prometheus.Registry) umami.Backend {
	return &prometheusBackend{
		registry: reg,
	}
}

func (p *prometheusBackend) Counter(opts umami.CounterOpts) umami.CounterBackend {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
	)
	p.registry.MustRegister(counter)
	return &prometheusCounterBackend{counter: counter}
}

func (p *prometheusBackend) CounterVec(opts umami.CounterVecOpts) umami.CounterVecBackend {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
		opts.Labels,
	)
	p.registry.MustRegister(counterVec)
	return &prometheusCounterVecBackend{counterVec: counterVec}
}

func (p *prometheusBackend) Gauge(opts umami.GaugeOpts) umami.GaugeBackend {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
	)
	p.registry.MustRegister(gauge)
	return &prometheusGaugeBackend{gauge: gauge}
}

func (p *prometheusBackend) GaugeVec(opts umami.GaugeVecOpts) umami.GaugeVecBackend {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
		opts.Labels,
	)
	p.registry.MustRegister(gaugeVec)
	return &prometheusGaugeVecBackend{gaugeVec: gaugeVec}
}

func (p *prometheusBackend) Histogram(opts umami.HistogramOpts) umami.HistogramBackend {
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    opts.Name,
			Help:    opts.Help,
			Buckets: opts.Buckets,
		},
	)
	p.registry.MustRegister(histogram)
	return &prometheusHistogramBackend{histogram: histogram}
}

func (p *prometheusBackend) HistogramVec(opts umami.HistogramVecOpts) umami.HistogramVecBackend {
	histogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    opts.Name,
			Help:    opts.Help,
			Buckets: opts.Buckets,
		},
		opts.Labels,
	)
	p.registry.MustRegister(histogramVec)
	return &prometheusHistogramVecBackend{histogramVec: histogramVec}
}

func (p *prometheusBackend) Summary(opts umami.SummaryOpts) umami.SummaryBackend {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		},
	)
	p.registry.MustRegister(summary)
	return &prometheusSummaryBackend{summary: summary}
}

func (p *prometheusBackend) SummaryVec(opts umami.SummaryVecOpts) umami.SummaryVecBackend {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		},
		opts.Labels,
	)
	p.registry.MustRegister(summaryVec)
	return &prometheusSummaryVecBackend{summaryVec: summaryVec}
}

func (p *prometheusBackend) Name() string {
	return PrometheusBackendName
}

var _prometheusBackend umami.Backend = (*prometheusBackend)(nil)
