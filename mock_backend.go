package umami

import (
	"fmt"
	"strings"
)

// mockBackend implements Backend interface for testing
type mockBackend struct {
	name string
}

// NewMockBackend creates a new mock backend for testing
func NewMockBackend() Backend {
	return &mockBackend{
		name: "mock",
	}
}

// GetMockBackend returns the mock backend for testing access to internals
func (m *mockBackend) AsMock() *mockBackend {
	return m
}

func (m *mockBackend) Name() string {
	return m.name
}

func (m *mockBackend) Counter(opts CounterOpts) CounterAdapter {
	return &mockCounterAdapter{
		name: opts.Name,
	}
}

func (m *mockBackend) CounterVec(opts CounterVecOpts) CounterVecAdapter {
	return &mockCounterVecAdapter{
		name:   opts.Name,
		counts: make(map[string]float64),
	}
}

func (m *mockBackend) Gauge(opts GaugeOpts) GaugeAdapter {
	return &mockGaugeAdapter{
		name: opts.Name,
	}
}

func (m *mockBackend) GaugeVec(opts GaugeVecOpts) GaugeVecAdapter {
	return &mockGaugeVecAdapter{
		name:   opts.Name,
		values: make(map[string]float64),
	}
}

func (m *mockBackend) Histogram(opts HistogramOpts) HistogramAdapter {
	return &mockHistogramAdapter{
		name: opts.Name,
	}
}

func (m *mockBackend) HistogramVec(opts HistogramVecOpts) HistogramVecAdapter {
	return &mockHistogramVecAdapter{
		name:         opts.Name,
		observations: make(map[string][]float64),
	}
}

func (m *mockBackend) Summary(opts SummaryOpts) SummaryAdapter {
	return &mockSummaryAdapter{
		name: opts.Name,
	}
}

func (m *mockBackend) SummaryVec(opts SummaryVecOpts) SummaryVecAdapater {
	return &mockSummaryVecAdapter{
		name:         opts.Name,
		observations: make(map[string][]float64),
	}
}

// Counter adapter
type mockCounterAdapter struct {
	name  string
	count float64
}

func (m *mockCounterAdapter) Inc() error {
	m.count++
	return nil
}

func (m *mockCounterAdapter) Add(value float64) error {
	m.count += value
	return nil
}

func (m *mockCounterAdapter) GetCount() float64 {
	return m.count
}

// CounterVec adapter
type mockCounterVecAdapter struct {
	name   string
	counts map[string]float64
}

func (m *mockCounterVecAdapter) Inc(labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.counts[key]++
	return nil
}

func (m *mockCounterVecAdapter) Add(value float64, labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.counts[key] += value
	return nil
}

func (m *mockCounterVecAdapter) GetCount(labels VecLabels) float64 {
	key := m.labelsToKey(labels)
	return m.counts[key]
}

func (m *mockCounterVecAdapter) labelsToKey(labels VecLabels) string {
	var parts []string
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}

// Gauge adapter
type mockGaugeAdapter struct {
	name  string
	value float64
}

func (m *mockGaugeAdapter) Set(value float64) error {
	m.value = value
	return nil
}

func (m *mockGaugeAdapter) Inc() error {
	m.value++
	return nil
}

func (m *mockGaugeAdapter) Dec() error {
	m.value--
	return nil
}

func (m *mockGaugeAdapter) Add(value float64) error {
	m.value += value
	return nil
}

func (m *mockGaugeAdapter) GetValue() float64 {
	return m.value
}

// GaugeVec adapter
type mockGaugeVecAdapter struct {
	name   string
	values map[string]float64
}

func (m *mockGaugeVecAdapter) Set(value float64, labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.values[key] = value
	return nil
}

func (m *mockGaugeVecAdapter) Inc(labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.values[key]++
	return nil
}

func (m *mockGaugeVecAdapter) Dec(labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.values[key]--
	return nil
}

func (m *mockGaugeVecAdapter) Add(value float64, labels VecLabels) error {
	key := m.labelsToKey(labels)
	m.values[key] += value
	return nil
}

func (m *mockGaugeVecAdapter) GetValue(labels VecLabels) float64 {
	key := m.labelsToKey(labels)
	return m.values[key]
}

func (m *mockGaugeVecAdapter) labelsToKey(labels VecLabels) string {
	var parts []string
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}

// Histogram adapter
type mockHistogramAdapter struct {
	name         string
	observations []float64
}

func (m *mockHistogramAdapter) Observe(value float64) error {
	m.observations = append(m.observations, value)
	return nil
}

func (m *mockHistogramAdapter) GetObservations() []float64 {
	return m.observations
}

func (m *mockHistogramAdapter) GetObservationCount() int {
	return len(m.observations)
}

// HistogramVec adapter
type mockHistogramVecAdapter struct {
	name         string
	observations map[string][]float64
}

func (m *mockHistogramVecAdapter) Observe(value float64, labels VecLabels) error {
	key := m.labelsToKey(labels)
	if m.observations[key] == nil {
		m.observations[key] = make([]float64, 0)
	}
	m.observations[key] = append(m.observations[key], value)
	return nil
}

func (m *mockHistogramVecAdapter) GetObservations(labels VecLabels) []float64 {
	key := m.labelsToKey(labels)
	return m.observations[key]
}

func (m *mockHistogramVecAdapter) GetObservationCount(labels VecLabels) int {
	key := m.labelsToKey(labels)
	return len(m.observations[key])
}

func (m *mockHistogramVecAdapter) labelsToKey(labels VecLabels) string {
	var parts []string
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}

// Summary adapter
type mockSummaryAdapter struct {
	name         string
	observations []float64
}

func (m *mockSummaryAdapter) Observe(value float64) error {
	m.observations = append(m.observations, value)
	return nil
}

func (m *mockSummaryAdapter) Quantile(q float64) (float64, error) {
	if len(m.observations) == 0 {
		return 0, nil
	}
	// Simple quantile calculation for testing
	index := int(q * float64(len(m.observations)))
	if index >= len(m.observations) {
		index = len(m.observations) - 1
	}
	return m.observations[index], nil
}

func (m *mockSummaryAdapter) GetObservations() []float64 {
	return m.observations
}

// SummaryVec adapter
type mockSummaryVecAdapter struct {
	name         string
	observations map[string][]float64
}

func (m *mockSummaryVecAdapter) Observe(value float64, labels VecLabels) error {
	key := m.labelsToKey(labels)
	if m.observations[key] == nil {
		m.observations[key] = make([]float64, 0)
	}
	m.observations[key] = append(m.observations[key], value)
	return nil
}

func (m *mockSummaryVecAdapter) Quantile(q float64, labels VecLabels) (float64, error) {
	key := m.labelsToKey(labels)
	obs := m.observations[key]
	if len(obs) == 0 {
		return 0, nil
	}
	// Simple quantile calculation for testing
	index := int(q * float64(len(obs)))
	if index >= len(obs) {
		index = len(obs) - 1
	}
	return obs[index], nil
}

func (m *mockSummaryVecAdapter) GetObservations(labels VecLabels) []float64 {
	key := m.labelsToKey(labels)
	return m.observations[key]
}

func (m *mockSummaryVecAdapter) labelsToKey(labels VecLabels) string {
	var parts []string
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ",")
}
