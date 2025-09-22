package umami_prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/SimonDaKappa/go-umami"
)

type prCounterAdapter struct {
	internal prometheus.Counter
}

func (pca *prCounterAdapter) Inc() error {
	pca.internal.Inc()
	return nil
}

func (pca *prCounterAdapter) Add(value float64) error {
	pca.internal.Add(value)
	return nil
}

type prCounterVecAdapter struct {
	internal *prometheus.CounterVec
}

func (pcva *prCounterVecAdapter) Inc(labels umami.VecLabels) error {
	pcva.internal.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (pcva *prCounterVecAdapter) Add(value float64, labels umami.VecLabels) error {
	pcva.internal.With(prometheus.Labels(labels)).Add(value)
	return nil
}

type prGaugeAdapter struct {
	internal prometheus.Gauge
}

func (pga *prGaugeAdapter) Set(value float64) error {
	pga.internal.Set(value)
	return nil
}

func (pga *prGaugeAdapter) Add(value float64) error {
	pga.internal.Add(value)
	return nil
}

func (pga *prGaugeAdapter) Inc() error {
	pga.internal.Inc()
	return nil
}

func (pga *prGaugeAdapter) Dec() error {
	pga.internal.Dec()
	return nil
}

type prGaugeVecAdapter struct {
	internal *prometheus.GaugeVec
}

func (pgva *prGaugeVecAdapter) Set(value float64, labels umami.VecLabels) error {
	pgva.internal.With(prometheus.Labels(labels)).Set(value)
	return nil
}

func (pgva *prGaugeVecAdapter) Add(value float64, labels umami.VecLabels) error {
	pgva.internal.With(prometheus.Labels(labels)).Add(value)
	return nil
}

func (pgva *prGaugeVecAdapter) Inc(labels umami.VecLabels) error {
	pgva.internal.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (pgva *prGaugeVecAdapter) Dec(labels umami.VecLabels) error {
	pgva.internal.With(prometheus.Labels(labels)).Dec()
	return nil
}

type prHistogramAdapter struct {
	internal prometheus.Histogram
}

func (pha *prHistogramAdapter) Observe(value float64) error {
	pha.internal.Observe(value)
	return nil
}

type prHistogramVecAdapter struct {
	internal *prometheus.HistogramVec
}

func (phva *prHistogramVecAdapter) Observe(value float64, labels umami.VecLabels) error {
	phva.internal.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

type prSummaryAdapter struct {
	internal prometheus.Summary
}

func (psa *prSummaryAdapter) Observe(value float64) error {
	psa.internal.Observe(value)
	return nil
}

func (psa *prSummaryAdapter) Quantile(q float64) (float64, error) {
	mfs := make(chan prometheus.Metric, 1)
	psa.internal.Collect(mfs)
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

type prSummaryVecAdapter struct {
	internal *prometheus.SummaryVec
}

func (m *prSummaryVecAdapter) Observe(value float64, labels umami.VecLabels) error {
	m.internal.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (m *prSummaryVecAdapter) Quantile(q float64, labels umami.VecLabels) (float64, error) {
	curried, err := m.internal.CurryWith(prometheus.Labels(labels))
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
	_pCounterBackend      umami.CounterAdapter      = (*prCounterAdapter)(nil)
	_pCounterVecBackend   umami.CounterVecAdapter   = (*prCounterVecAdapter)(nil)
	_pGaugeBackend        umami.GaugeAdapter        = (*prGaugeAdapter)(nil)
	_pGaugeVecBackend     umami.GaugeVecAdapter     = (*prGaugeVecAdapter)(nil)
	_pHistogramBackend    umami.HistogramAdapter    = (*prHistogramAdapter)(nil)
	_pHistogramVecBackend umami.HistogramVecAdapter = (*prHistogramVecAdapter)(nil)
)
