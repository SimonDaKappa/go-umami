package umami_backend

// Backend interface allows plugging in different 3rd party monitoring services
type Backend interface {
	// Counter creates a counter metric
	Counter(name string, labels ...string) CounterBackend

	// Gauge creates a gauge metric
	Gauge(name string, labels ...string) GaugeBackend

	// Histogram creates a histogram metric
	Histogram(name string, buckets []float64, labels ...string) HistogramBackend

	// Summary creates a summary metric
	Summary(name string, objectives map[float64]float64, labels ...string) SummaryBackend
}

// CounterBackend defines the interface for counter metrics that concrete
// backend abstraction implementations must satisfy.
type CounterBackend interface {
	Inc() error
	Add(value float64) error
}

// GaugeBackend defines the interface for gauge metrics that concrete
// backend abstraction implementations must satisfy.
type GaugeBackend interface {
	Set(value float64) error
	Inc() error
	Dec() error
	Add(value float64) error
}

// HistogramBackend defines the interface for histogram metrics that concrete
// backend abstraction implementations must satisfy.
type HistogramBackend interface {
	Observe(value float64) error
}

// SummaryBackend defines the interface for summary metrics that concrete
// backend abstraction implementations must satisfy.
type SummaryBackend interface {
	Observe(value float64) error
	Quantile(q float64) (float64, error)
}
