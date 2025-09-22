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

func (p *prometheusBackend) Counter(opts umami.CounterOpts) umami.CounterAdapter {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
	)
	p.registry.MustRegister(counter)
	return &prCounterAdapter{internal: counter}
}

func (p *prometheusBackend) CounterVec(opts umami.CounterVecOpts) umami.CounterVecAdapter {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
		opts.Labels,
	)
	p.registry.MustRegister(counterVec)
	return &prCounterVecAdapter{internal: counterVec}
}

func (p *prometheusBackend) Gauge(opts umami.GaugeOpts) umami.GaugeAdapter {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
	)
	p.registry.MustRegister(gauge)
	return &prGaugeAdapter{internal: gauge}
}

func (p *prometheusBackend) GaugeVec(opts umami.GaugeVecOpts) umami.GaugeVecAdapter {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: opts.Name,
			Help: opts.Help,
		},
		opts.Labels,
	)
	p.registry.MustRegister(gaugeVec)
	return &prGaugeVecAdapter{internal: gaugeVec}
}

func (p *prometheusBackend) Histogram(opts umami.HistogramOpts) umami.HistogramAdapter {
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    opts.Name,
			Help:    opts.Help,
			Buckets: opts.Buckets,
		},
	)
	p.registry.MustRegister(histogram)
	return &prHistogramAdapter{internal: histogram}
}

func (p *prometheusBackend) HistogramVec(opts umami.HistogramVecOpts) umami.HistogramVecAdapter {
	histogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    opts.Name,
			Help:    opts.Help,
			Buckets: opts.Buckets,
		},
		opts.Labels,
	)
	p.registry.MustRegister(histogramVec)
	return &prHistogramVecAdapter{internal: histogramVec}
}

func (p *prometheusBackend) Summary(opts umami.SummaryOpts) umami.SummaryAdapter {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		},
	)
	p.registry.MustRegister(summary)
	return &prSummaryAdapter{internal: summary}
}

func (p *prometheusBackend) SummaryVec(opts umami.SummaryVecOpts) umami.SummaryVecAdapater {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		},
		opts.Labels,
	)
	p.registry.MustRegister(summaryVec)
	return &prSummaryVecAdapter{internal: summaryVec}
}

func (p *prometheusBackend) Name() string {
	return PrometheusBackendName
}

var __ctc_prometheusBackend umami.Backend = (*prometheusBackend)(nil)
