package integrations

import (
	"fmt"
	"github.com/SimonDaKappa/"
)

// Mock backend for demonstration
type prometheusBackend struct{}

func NewPrometheusBackend() Backend {
	return &prometheusBackend{}
}

func (p *prometheusBackend) Counter(name string, labels ...string) CounterBackend {
	return &mockCounterBackend{name: name}
}

func (p *prometheusBackend) Gauge(name string, labels ...string) GaugeBackend {
	return &mockGaugeBackend{name: name}
}

func (p *prometheusBackend) Histogram(name string, buckets []float64, labels ...string) HistogramBackend {
	return &mockHistogramBackend{name: name}
}

func (p *prometheusBackend) Summary(name string, objectives map[float64]float64, labels ...string) SummaryBackend {
	return &mockSummaryBackend{name: name}
}

// Mock implementations
type mockCounterBackend struct{ name string }

func (m *mockCounterBackend) Inc() error { fmt.Printf("Counter %s incremented\n", m.name); return nil }
func (m *mockCounterBackend) Add(value float64) error {
	fmt.Printf("Counter %s added %f\n", m.name, value)
	return nil
}

type mockGaugeBackend struct{ name string }

func (m *mockGaugeBackend) Set(value float64) error {
	fmt.Printf("Gauge %s set to %f\n", m.name, value)
	return nil
}
func (m *mockGaugeBackend) Inc() error { fmt.Printf("Gauge %s incremented\n", m.name); return nil }
func (m *mockGaugeBackend) Dec() error { fmt.Printf("Gauge %s decremented\n", m.name); return nil }
func (m *mockGaugeBackend) Add(value float64) error {
	fmt.Printf("Gauge %s added %f\n", m.name, value)
	return nil
}

type mockHistogramBackend struct{ name string }

func (m *mockHistogramBackend) Observe(value float64) error {
	fmt.Printf("Histogram %s observed %f\n", m.name, value)
	return nil
}

type mockSummaryBackend struct{ name string }

func (m *mockSummaryBackend) Observe(value float64) error {
	fmt.Printf("Summary %s observed %f\n", m.name, value)
	return nil
}
func (m *mockSummaryBackend) Quantile(q float64) (float64, error) { return q * 100, nil }
