package umami

import (
	"fmt"
	"time"
)

// Example usage of the redesigned metrics system

func ExampleUsage() {
	// 1. Create a backend (could be Prometheus, DataDog, StatsD, etc.)
	backend := NewPrometheusBackend() // This would be implemented separately

	// 2. Create metrics manager
	manager := NewManager(backend)

	// 3. Configure global settings
	manager.SetGlobalLevel(LevelImportant)
	manager.SetGlobalMask(MaskProduction)

	// 4. Create metric groups for different components
	webGroup := manager.Group("web")
	_ = manager.Group("database") // Could be used for DB metrics
	_ = manager.Group("pipeline") // Could be used for pipeline metrics

	// 5. Create metrics using the factory
	webFactory := webGroup.Factory()

	// Critical metrics - always enabled in production
	requestCounter := webFactory.Counter("requests_total", LevelCritical, MaskCounters, "method", "status")
	requestLatency := webFactory.Histogram("request_duration_seconds", LevelCritical, MaskLatency,
		[]float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}, "method")

	// Important operational metrics
	dbConnPool := webFactory.Pool("db_connections", LevelImportant, MaskConnections)

	// Debug metrics - often disabled in production
	perUserMetrics := webFactory.Counter("user_actions", LevelDebug, MaskPerUser, "user_id", "action")

	// 6. Usage in application code
	ctx := webGroup.Context()

	// Simulate handling a web request
	handleRequest(ctx, requestCounter, requestLatency, dbConnPool, perUserMetrics)

	// 7. Runtime configuration changes
	// Disable debug metrics in production
	webGroup.SetLevel(LevelImportant) // This will disable LevelDebug metrics

	// Or disable specific metric types
	webGroup.SetMask(MaskEssential) // Only keep counters, latency, and errors
}

func handleRequest(ctx Context, requestCounter Counter, requestLatency Histogram,
	dbPool Pool, userMetrics Counter) {

	// Early return pattern - if metrics are disabled, these calls return immediately
	requestCounter.Inc(ctx)

	// Time the request
	timer := requestLatency.Time(ctx, func() error {
		// Simulate request processing
		time.Sleep(10 * time.Millisecond)

		// Record database pool usage
		dbPool.Acquired(ctx)
		defer dbPool.Released(ctx)

		// This might be disabled in production (LevelDebug)
		userMetrics.Add(ctx, 1)

		return nil
	})

	fmt.Printf("Request processed, timer error: %v\n", timer)
}

// Configuration example
func ConfigurationExample() {
	backend := NewPrometheusBackend()
	manager := NewManager(backend)

	// Production configuration
	manager.SetGlobalLevel(LevelImportant)
	manager.SetGlobalMask(MaskProduction)

	// Development configuration
	// manager.SetGlobalLevel(LevelVerbose)
	// manager.SetGlobalMask(MaskAll)

	// Per-group configuration
	webGroup := manager.Group("web")
	webGroup.SetLevel(LevelCritical) // Only critical web metrics

	debugGroup := manager.Group("debug")
	debugGroup.SetLevel(LevelVerbose) // All debug metrics
	debugGroup.SetMask(MaskInternal | MaskDetailed)
}

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
