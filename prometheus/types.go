package umami_prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/SimonDaKappa/go-umami"
)

type prometheusCounterBackend struct {
	counter prometheus.Counter
}

func (m *prometheusCounterBackend) Inc() error {
	m.counter.Inc()
	return nil
}

func (m *prometheusCounterBackend) Add(value float64) error {
	m.counter.Add(value)
	return nil
}

type prometheusCounterVecBackend struct {
	counterVec *prometheus.CounterVec
}

func (m *prometheusCounterVecBackend) Inc(labels umami.VecLabels) error {
	m.counterVec.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (m *prometheusCounterVecBackend) Add(value float64, labels umami.VecLabels) error {
	m.counterVec.With(prometheus.Labels(labels)).Add(value)
	return nil
}

type prometheusGaugeBackend struct {
	gauge prometheus.Gauge
}

func (m *prometheusGaugeBackend) Set(value float64) error {
	m.gauge.Set(value)
	return nil
}

func (m *prometheusGaugeBackend) Add(value float64) error {
	m.gauge.Add(value)
	return nil
}

func (m *prometheusGaugeBackend) Inc() error {
	m.gauge.Inc()
	return nil
}

func (m *prometheusGaugeBackend) Dec() error {
	m.gauge.Dec()
	return nil
}

type prometheusGaugeVecBackend struct {
	gaugeVec *prometheus.GaugeVec
}

func (m *prometheusGaugeVecBackend) Set(value float64, labels umami.VecLabels) error {
	m.gaugeVec.With(prometheus.Labels(labels)).Set(value)
	return nil
}

func (m *prometheusGaugeVecBackend) Add(value float64, labels umami.VecLabels) error {
	m.gaugeVec.With(prometheus.Labels(labels)).Add(value)
	return nil
}

func (m *prometheusGaugeVecBackend) Inc(labels umami.VecLabels) error {
	m.gaugeVec.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (m *prometheusGaugeVecBackend) Dec(labels umami.VecLabels) error {
	m.gaugeVec.With(prometheus.Labels(labels)).Dec()
	return nil
}

type prometheusHistogramBackend struct {
	histogram prometheus.Histogram
}

func (m *prometheusHistogramBackend) Observe(value float64) error {
	m.histogram.Observe(value)
	return nil
}

type prometheusHistogramVecBackend struct {
	histogramVec *prometheus.HistogramVec
}

func (m *prometheusHistogramVecBackend) Observe(value float64, labels umami.VecLabels) error {
	m.histogramVec.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

type prometheusSummaryBackend struct {
	summary prometheus.Summary
}

func (m *prometheusSummaryBackend) Observe(value float64) error {
	m.summary.Observe(value)
	return nil
}

func (m *prometheusSummaryBackend) Quantile(q float64) (float64, error) {
	mfs := make(chan prometheus.Metric, 1)
	m.summary.Collect(mfs)
	close(mfs)
	for metric := range mfs {
		dto := &dto.Metric{}
		if err := metric.Write(dto); err != nil {
			return 0, err
		}
		for _, quantile := range dto.Summary.Quantile {
			if quantile.GetQuantile() == q {
				return quantile.GetValue(), nil
			}
		}
	}
	return 0, fmt.Errorf("Quantile %f not found", q)
}

type prometheusSummaryVecBackend struct {
	summaryVec *prometheus.SummaryVec
}

func (m *prometheusSummaryVecBackend) Observe(value float64, labels umami.VecLabels) error {
	m.summaryVec.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (m *prometheusSummaryVecBackend) Quantile(q float64, labels umami.VecLabels) (float64, error) {
	curried, err := m.summaryVec.CurryWith(prometheus.Labels(labels))
	if err != nil {
		return 0, err
	}

	mfs := make(chan prometheus.Metric, 1)
	curried.Collect(mfs)
	close(mfs)

	for metric := range mfs {
		dto := &dto.Metric{}
		if err := metric.Write(dto); err != nil {
			return 0, err
		}
		for _, quantile := range dto.Summary.Quantile {
			if quantile.GetQuantile() == q {
				return quantile.GetValue(), nil
			}
		}
	}
	return 0, fmt.Errorf("Quantile %f not found", q)
}

// Sanity checks for interface implementation
var (
	_pCounterBackend      umami.CounterBackend      = (*prometheusCounterBackend)(nil)
	_pCounterVecBackend   umami.CounterVecBackend   = (*prometheusCounterVecBackend)(nil)
	_pGaugeBackend        umami.GaugeBackend        = (*prometheusGaugeBackend)(nil)
	_pGaugeVecBackend     umami.GaugeVecBackend     = (*prometheusGaugeVecBackend)(nil)
	_pHistogramBackend    umami.HistogramBackend    = (*prometheusHistogramBackend)(nil)
	_pHistogramVecBackend umami.HistogramVecBackend = (*prometheusHistogramVecBackend)(nil)
)
