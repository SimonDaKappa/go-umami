package umami

import "time"

// Context allows checking if metrics are enabled without coupling to specific implementations
type Context interface {
	// Enabled returns true if metrics at this level should be processed
	Enabled(level Level) bool

	// EnabledMask returns true if metrics with this mask should be processed
	EnabledMask(mask MetricMask) bool

	// WithLevel returns a new context with the specified level
	WithLevel(level Level) Context

	// WithMask returns a new context with the specified mask
	WithMask(mask MetricMask) Context
}

type Counter interface {
	// Inc increments the counter. Returns immediately if metric is disabled.
	Inc(ctx Context) error

	// Add adds the given value to the counter. Returns immediately if metric is disabled.
	Add(ctx Context, value float64) error
}

type Gauge interface {
	// Set sets the gauge to the given value. Returns immediately if metric is disabled.
	Set(ctx Context, value float64) error

	// Inc increments the gauge. Returns immediately if metric is disabled.
	Inc(ctx Context) error

	// Dec decrements the gauge. Returns immediately if metric is disabled.
	Dec(ctx Context) error

	// Add adds the given value to the gauge. Returns immediately if metric is disabled.
	Add(ctx Context, value float64) error
}

type Histogram interface {
	// Observe adds an observation to the histogram. Returns immediately if metric is disabled.
	Observe(ctx Context, value float64) error

	// Time times the execution of the function. Returns immediately if metric is disabled.
	Time(ctx Context, fn func() error) error
}

type Summary interface {
	// Observe adds an observation to the summary. Returns immediately if metric is disabled.
	Observe(ctx Context, value float64) error

	// Quantile returns the value at the given quantile. Returns 0 if metric is disabled.
	Quantile(ctx Context, q float64) (float64, error)
}

// Specialized metric interfaces for common backend patterns
type Timer interface {
	// Start returns a function that should be called when the operation completes
	// Returns a no-op function if metric is disabled
	Start(ctx Context) func()

	// Record records a duration. Returns immediately if metric is disabled.
	Record(ctx Context, duration time.Duration) error
}

type Cache interface {
	// Hit records a cache hit. Returns immediately if metric is disabled.
	Hit(ctx Context) error

	// Miss records a cache miss. Returns immediately if metric is disabled.
	Miss(ctx Context) error

	// SetSize sets the current cache size. Returns immediately if metric is disabled.
	SetSize(ctx Context, bytes int64) error
}

type Pool interface {
	// SetActive sets the number of active connections. Returns immediately if metric is disabled.
	SetActive(ctx Context, count int) error

	// SetIdle sets the number of idle connections. Returns immediately if metric is disabled.
	SetIdle(ctx Context, count int) error

	// Acquired records a connection acquisition. Returns immediately if metric is disabled.
	Acquired(ctx Context) error

	// Released records a connection release. Returns immediately if metric is disabled.
	Released(ctx Context) error
}

type CircuitBreaker interface {
	// SetState sets the circuit breaker state. Returns immediately if metric is disabled.
	SetState(ctx Context, state string) error

	// Success records a successful operation. Returns immediately if metric is disabled.
	Success(ctx Context) error

	// Failure records a failed operation. Returns immediately if metric is disabled.
	Failure(ctx Context) error
}

type Queue interface {
	// SetDepth sets the current queue depth. Returns immediately if metric is disabled.
	SetDepth(ctx Context, depth int) error

	// Enqueued records an item being enqueued. Returns immediately if metric is disabled.
	Enqueued(ctx Context) error

	// Dequeued records an item being dequeued. Returns immediately if metric is disabled.
	Dequeued(ctx Context) error

	// SetWaitTime records how long items wait in the queue. Returns immediately if metric is disabled.
	SetWaitTime(ctx Context, duration time.Duration) error
}

// Factory creates metrics with the appropriate level and mask
type Factory interface {
	// Counter creates a counter with the given level and mask
	Counter(name string, level Level, mask MetricMask, labels ...string) Counter

	// Gauge creates a gauge with the given level and mask
	Gauge(name string, level Level, mask MetricMask, labels ...string) Gauge

	// Histogram creates a histogram with the given level and mask
	Histogram(name string, level Level, mask MetricMask, buckets []float64, labels ...string) Histogram

	// Summary creates a summary with the given level and mask
	Summary(name string, level Level, mask MetricMask, objectives map[float64]float64, labels ...string) Summary

	// Timer creates a timer with the given level and mask
	Timer(name string, level Level, mask MetricMask, labels ...string) Timer

	// Cache creates cache metrics with the given level and mask
	Cache(name string, level Level, mask MetricMask, labels ...string) Cache

	// Pool creates pool metrics with the given level and mask
	Pool(name string, level Level, mask MetricMask, labels ...string) Pool

	// CircuitBreaker creates circuit breaker metrics with the given level and mask
	CircuitBreaker(name string, level Level, mask MetricMask, labels ...string) CircuitBreaker

	// Queue creates queue metrics with the given level and mask
	Queue(name string, level Level, mask MetricMask, labels ...string) Queue
}

// Group represents a logical grouping of metrics (e.g., "web", "database", "pipeline")
type Group interface {
	// Factory returns a factory for creating metrics in this group
	Factory() Factory

	// SetLevel sets the level for this group
	SetLevel(level Level)

	// SetMask sets the mask for this group
	SetMask(mask MetricMask)

	// Context returns a context for this group
	Context() Context
}

// Manager is the top-level metrics manager
type Manager interface {
	// Group returns or creates a metric group
	Group(name string) Group

	// SetGlobalLevel sets the global metrics level
	SetGlobalLevel(level Level)

	// SetGlobalMask sets the global metrics mask
	SetGlobalMask(mask MetricMask)

	// GlobalContext returns the global metrics context
	GlobalContext() Context
}
