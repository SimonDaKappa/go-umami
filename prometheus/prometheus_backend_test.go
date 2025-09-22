package umami_prometheus

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/SimonDaKappa/go-umami"
)

func TestManagerIntegrationWithPrometheusBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	manager := umami.NewRegistry(backend)
	group := manager.Group("web")
	factory := group.Factory()

	counter := factory.Counter(umami.CounterOpts{Name: "integration_counter", Help: "integration counter"}, 0)
	ctx := umami.NewContext(0)
	if err := counter.Inc(ctx); err != nil {
		t.Errorf("Manager integration Counter Inc failed: %v", err)
	}
	if err := counter.Add(ctx, 2); err != nil {
		t.Errorf("Manager integration Counter Add failed: %v", err)
	}

	gauge := factory.Gauge(umami.GaugeOpts{Name: "integration_gauge", Help: "integration gauge"}, 0)
	if err := gauge.Set(ctx, 10); err != nil {
		t.Errorf("Manager integration Gauge Set failed: %v", err)
	}
	if err := gauge.Inc(ctx); err != nil {
		t.Errorf("Manager integration Gauge Inc failed: %v", err)
	}
	if err := gauge.Dec(ctx); err != nil {
		t.Errorf("Manager integration Gauge Dec failed: %v", err)
	}
	if err := gauge.Add(ctx, 5); err != nil {
		t.Errorf("Manager integration Gauge Add failed: %v", err)
	}

	hist := factory.Histogram(umami.HistogramOpts{Name: "integration_histogram", Help: "integration histogram", Buckets: []float64{0.1, 1, 10}}, 0)
	if err := hist.Observe(ctx, 0.7); err != nil {
		t.Errorf("Manager integration Histogram Observe failed: %v", err)
	}

	summary := factory.Summary(umami.SummaryOpts{Name: "integration_summary", Help: "integration summary", Objectives: map[float64]float64{0.5: 0.05}}, 0)
	if err := summary.Observe(ctx, 1.5); err != nil {
		t.Errorf("Manager integration Summary Observe failed: %v", err)
	}
	_, _ = summary.Quantile(ctx, 0.5)
}

func TestPrometheusCounterBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	counter := backend.Counter(umami.CounterOpts{Name: "test_counter", Help: "test counter"})
	if err := counter.Inc(); err != nil {
		t.Errorf("Counter Inc failed: %v", err)
	}
	if err := counter.Add(5); err != nil {
		t.Errorf("Counter Add failed: %v", err)
	}
}

func TestPrometheusGaugeBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	gauge := backend.Gauge(umami.GaugeOpts{Name: "test_gauge", Help: "test gauge"})
	if err := gauge.Set(42); err != nil {
		t.Errorf("Gauge Set failed: %v", err)
	}
	if err := gauge.Inc(); err != nil {
		t.Errorf("Gauge Inc failed: %v", err)
	}
	if err := gauge.Dec(); err != nil {
		t.Errorf("Gauge Dec failed: %v", err)
	}
	if err := gauge.Add(3); err != nil {
		t.Errorf("Gauge Add failed: %v", err)
	}
}

func TestPrometheusHistogramBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	hist := backend.Histogram(umami.HistogramOpts{Name: "test_histogram", Help: "test histogram", Buckets: []float64{0.1, 1, 10}})
	if err := hist.Observe(0.5); err != nil {
		t.Errorf("Histogram Observe failed: %v", err)
	}
}

func TestPrometheusSummaryBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	summary := backend.Summary(umami.SummaryOpts{Name: "test_summary", Help: "test summary", Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01}})
	if err := summary.Observe(1.0); err != nil {
		t.Errorf("Summary Observe failed: %v", err)
	}
	_, _ = summary.Quantile(0.5)
}

func getMetricValue(t *testing.T, reg *prometheus.Registry, name string, labelPairs map[string]string) float64 {
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	for _, mf := range mfs {
		if mf.GetName() == name {
			for _, m := range mf.GetMetric() {
				match := true
				for k, v := range labelPairs {
					found := false
					for _, lp := range m.GetLabel() {
						if lp.GetName() == k && lp.GetValue() == v {
							found = true
							break
						}
					}
					if !found {
						match = false
						break
					}
				}
				if match {
					if m.Gauge != nil {
						return m.Gauge.GetValue()
					}
					if m.Counter != nil {
						return m.Counter.GetValue()
					}
					if m.Histogram != nil {
						return m.Histogram.GetSampleSum()
					}
					if m.Summary != nil {
						return m.Summary.GetSampleSum()
					}
				}
			}
		}
	}
	t.Fatalf("metric %s with labels %+v not found", name, labelPairs)
	return 0
}

func TestPrometheusCounterVecBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	counterVec := backend.CounterVec(umami.CounterVecOpts{Name: "test_counter_vec", Help: "test counter vec", Labels: []string{"foo", "bar"}})
	labels := umami.VecLabels{"foo": "a", "bar": "b"}
	if err := counterVec.Inc(labels); err != nil {
		t.Errorf("CounterVec Inc failed: %v", err)
	}
	if err := counterVec.Add(3, labels); err != nil {
		t.Errorf("CounterVec Add failed: %v", err)
	}
	val := getMetricValue(t, reg, "test_counter_vec", map[string]string{"foo": "a", "bar": "b"})
	if val != 4 {
		t.Errorf("CounterVec value = %v, want 4", val)
	}
}

func TestPrometheusGaugeVecBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	gaugeVec := backend.GaugeVec(umami.GaugeVecOpts{Name: "test_gauge_vec", Help: "test gauge vec", Labels: []string{"foo"}})
	labels := umami.VecLabels{"foo": "bar"}
	if err := gaugeVec.Set(10, labels); err != nil {
		t.Errorf("GaugeVec Set failed: %v", err)
	}
	if err := gaugeVec.Inc(labels); err != nil {
		t.Errorf("GaugeVec Inc failed: %v", err)
	}
	if err := gaugeVec.Dec(labels); err != nil {
		t.Errorf("GaugeVec Dec failed: %v", err)
	}
	if err := gaugeVec.Add(5, labels); err != nil {
		t.Errorf("GaugeVec Add failed: %v", err)
	}
	val := getMetricValue(t, reg, "test_gauge_vec", map[string]string{"foo": "bar"})
	if val != 15 {
		t.Errorf("GaugeVec value = %v, want 15", val)
	}
}

func TestPrometheusHistogramVecBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	histVec := backend.HistogramVec(umami.HistogramVecOpts{Name: "test_histogram_vec", Help: "test histogram vec", Buckets: []float64{0.1, 1, 10}, Labels: []string{"foo"}})
	labels := umami.VecLabels{"foo": "bar"}
	if err := histVec.Observe(0.5, labels); err != nil {
		t.Errorf("HistogramVec Observe failed: %v", err)
	}
	val := getMetricValue(t, reg, "test_histogram_vec", map[string]string{"foo": "bar"})
	if val != 0.5 {
		t.Errorf("HistogramVec value = %v, want 0.5", val)
	}
}

func TestPrometheusSummaryVecBackend(t *testing.T) {
	reg := prometheus.NewRegistry()
	backend := NewPrometheusBackend(reg)
	summaryVec := backend.SummaryVec(umami.SummaryVecOpts{Name: "test_summary_vec", Help: "test summary vec", Objectives: map[float64]float64{0.5: 0.05}, Labels: []string{"foo"}})
	labels := umami.VecLabels{"foo": "bar"}
	if err := summaryVec.Observe(2.0, labels); err != nil {
		t.Errorf("SummaryVec Observe failed: %v", err)
	}
	// Give Prometheus a moment to process
	time.Sleep(10 * time.Millisecond)
	val := getMetricValue(t, reg, "test_summary_vec", map[string]string{"foo": "bar"})
	if val != 2.0 {
		t.Errorf("SummaryVec value = %v, want 2.0", val)
	}
}
