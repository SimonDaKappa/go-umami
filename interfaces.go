package umami

import "time"

// Context allows checking if metrics are enabled without coupling to specific implementations
type Context interface {
	// Enabled returns true if metrics at this level should be processed
	Enabled(level Level) bool

	// EnabledMask returns true if metrics with this mask should be processed
	EnabledMask(mask Mask) bool

	// WithLevel returns a new context with the specified level
	WithLevel(level Level) Context

	// WithMask returns a new context with the specified mask
	WithMask(mask Mask) Context
}

type Counter interface {
	// Inc increments the counter. Noop if disabled.
	Inc(ctx Context) error

	// Add adds the given value to the counter. Noop if disabled.
	Add(ctx Context, value float64) error
}

type CounterVec interface {
	// Inc increments the counter for the given labels. Noop if disabled.
	Inc(ctx Context, labels VecLabels) error

	// Add adds the given value to the counter for the given labels. Noop if disabled.
	Add(ctx Context, value float64, labels VecLabels) error
}

type Gauge interface {
	// Set sets the gauge to the given value. Noop if disabled.
	Set(ctx Context, value float64) error

	// Inc increments the gauge. Noop if disabled.
	Inc(ctx Context) error

	// Dec decrements the gauge. Noop if disabled.
	Dec(ctx Context) error

	// Add adds the given value to the gauge. Noop if disabled.
	Add(ctx Context, value float64) error
}

type GaugeVec interface {
	// Set sets the gauge for the given labels to the given value. Noop if disabled.
	Set(ctx Context, value float64, labels VecLabels) error

	// Inc increments the gauge for the given labels. Noop if disabled.
	Inc(ctx Context, labels VecLabels) error

	// Dec decrements the gauge for the given labels. Noop if disabled.
	Dec(ctx Context, labels VecLabels) error

	// Add adds the given value to the gauge for the given labels. Noop if disabled.
	Add(ctx Context, value float64, labels VecLabels) error
}

type Histogram interface {
	// Observe adds an observation to the histogram. Noop if disabled.
	Observe(ctx Context, value float64) error

	// Time times the execution of the function. Noop if disabled.
	Time(ctx Context, fn func() error) error
}

type HistogramVec interface {
	// Observe adds an observation to the histogram for the given labels. Noop if disabled.
	Observe(ctx Context, value float64, labels VecLabels) error

	// Time times the execution of the function for the given labels. Noop if disabled.
	Time(ctx Context, fn func() error, labels VecLabels) error
}

type Summary interface {
	// Observe adds an observation to the summary. Noop if disabled.
	Observe(ctx Context, value float64) error

	// Quantile returns the value at the given quantile. Returns 0 if metric is disabled.
	Quantile(ctx Context, q float64) (float64, error)
}

type SummaryVec interface {
	// Observe adds an observation to the summary for the given labels. Noop if disabled.
	Observe(ctx Context, value float64, labels VecLabels) error

	// Quantile returns the value at the given quantile for the given labels. Returns 0 if metric is disabled.
	Quantile(ctx Context, q float64, labels VecLabels) (float64, error)
}

// Specialized metric interfaces for common backend patterns
type Timer interface {
	// Start returns a function that should be called when the operation completes
	// Returns a no-op function if metric is disabled
	Start(ctx Context) func()

	// Record records a duration. Noop if disabled.
	Record(ctx Context, duration time.Duration) error
}

type TimerVec interface {
	// Start returns a function that should be called when the operation completes
	// Returns a no-op function if metric is disabled
	Start(ctx Context, labels VecLabels) func()

	// Record records a duration. Noop if disabled.
	Record(ctx Context, duration time.Duration, labels VecLabels) error
}

type Cache interface {
	// Hit records a cache hit. Noop if disabled.
	Hit(ctx Context) error

	// Miss records a cache miss. Noop if disabled.
	Miss(ctx Context) error

	// SetSize sets the current cache size. Noop if disabled.
	SetSize(ctx Context, bytes int64) error
}

type CacheVec interface {
	// Hit records a cache hit for the given labels. Noop if disabled.
	Hit(ctx Context, labels VecLabels) error

	// Miss records a cache miss for the given labels. Noop if disabled.
	Miss(ctx Context, labels VecLabels) error

	// SetSize sets the current cache size for the given labels. Noop if disabled.
	SetSize(ctx Context, bytes int64, labels VecLabels) error
}

type Pool interface {
	// SetActive sets the number of active connections. Noop if disabled.
	SetActive(ctx Context, count int) error

	// SetIdle sets the number of idle connections. Noop if disabled.
	SetIdle(ctx Context, count int) error

	// Acquired records a connection acquisition. Noop if disabled.
	Acquired(ctx Context) error

	// Released records a connection release. Noop if disabled.
	Released(ctx Context) error
}

type PoolVec interface {
	// SetActive sets the number of active connections for the given labels. Noop if disabled.
	SetActive(ctx Context, count int, labels VecLabels) error

	// SetIdle sets the number of idle connections for the given labels. Noop if disabled.
	SetIdle(ctx Context, count int, labels VecLabels) error

	// Acquired records a connection acquisition for the given labels. Noop if disabled.
	Acquired(ctx Context, labels VecLabels) error

	// Released records a connection release for the given labels. Noop if disabled.
	Released(ctx Context, labels VecLabels) error
}

type CircuitBreaker interface {
	// SetState sets the circuit breaker state. Noop if disabled.
	SetState(ctx Context, state CircuitBreakerState) error

	// Success records a successful operation. Noop if disabled.
	Success(ctx Context) error

	// Failure records a failed operation. Noop if disabled.
	Failure(ctx Context) error
}

type CircuitBreakerVec interface {
	// SetState sets the circuit breaker state for the given labels. Noop if disabled.
	SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error

	// Success records a successful operation for the given labels. Noop if disabled.
	Success(ctx Context, labels VecLabels) error

	// Failure records a failed operation for the given labels. Noop if disabled.
	Failure(ctx Context, labels VecLabels) error
}

type Queue interface {
	// SetDepth sets the current queue depth. Noop if disabled.
	SetDepth(ctx Context, depth int) error

	// Enqueued records an item being enqueued. Noop if disabled.
	Enqueued(ctx Context) error

	// Dequeued records an item being dequeued. Noop if disabled.
	Dequeued(ctx Context) error

	// SetWaitTime records how long items wait in the queue. Noop if disabled.
	SetWaitTime(ctx Context, duration time.Duration) error
}

type QueueVec interface {
	// SetDepth sets the current queue depth for the given labels. Noop if disabled.
	SetDepth(ctx Context, depth int, labels VecLabels) error

	// Enqueued records an item being enqueued for the given labels. Noop if disabled.
	Enqueued(ctx Context, labels VecLabels) error

	// Dequeued records an item being dequeued for the given labels. Noop if disabled.
	Dequeued(ctx Context, labels VecLabels) error

	// SetWaitTime records how long items wait in the queue for the given labels. Noop if disabled.
	SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error
}

// Factory creates metrics with the appropriate level and mask
type Factory interface {
	// Counter creates a counter with the given level and mask
	Counter(opts CounterOpts, level Level, mask Mask) Counter

	// CounterVec creates a label-vectorized counter with the given level and mask
	CounterVec(opts CounterVecOpts, level Level, mask Mask) CounterVec

	// Gauge creates a gauge with the given level and mask
	Gauge(opts GaugeOpts, level Level, mask Mask) Gauge

	// GaugeVec creates a label-vectorized gauge with the given level and mask
	GaugeVec(opts GaugeVecOpts, level Level, mask Mask) GaugeVec

	// Histogram creates a histogram with the given level and mask
	Histogram(opts HistogramOpts, level Level, mask Mask) Histogram

	// HistogramVec creates a label-vectorized histogram with the given level and mask
	HistogramVec(opts HistogramVecOpts, level Level, mask Mask) HistogramVec

	// Summary creates a summary with the given level and mask
	Summary(opts SummaryOpts, level Level, mask Mask) Summary

	// SummaryVec creates a label-vectorized summary with the given level and mask
	SummaryVec(opts SummaryVecOpts, level Level, mask Mask) SummaryVec

	// Timer creates a timer with the given level and mask
	Timer(opts TimerOpts, level Level, mask Mask) Timer

	// TimerVec creates a label-vectorized timer with the given level and mask
	TimerVec(opts TimerVecOpts, level Level, mask Mask) TimerVec

	// Cache creates cache metrics with the given level and mask
	Cache(opts CacheOpts, level Level, mask Mask) Cache

	// CacheVec creates a label-vectorized cache with the given level and mask
	CacheVec(opts CacheVecOpts, level Level, mask Mask) CacheVec

	// Pool creates pool metrics with the given level and mask
	Pool(opts PoolOpts, level Level, mask Mask) Pool

	// PoolVec creates a label-vectorized pool with the given level and mask
	PoolVec(opts PoolVecOpts, level Level, mask Mask) PoolVec

	// CircuitBreaker creates circuit breaker metrics with the given level and mask
	CircuitBreaker(opts CircuitBreakerOpts, level Level, mask Mask) CircuitBreaker

	// CircuitBreakerVec creates a label-vectorized circuit breaker with the given level and mask
	CircuitBreakerVec(opts CircuitBreakerVecOpts, level Level, mask Mask) CircuitBreakerVec

	// Queue creates queue metrics with the given level and mask
	Queue(opts QueueOpts, level Level, mask Mask) Queue

	// QueueVec creates a label-vectorized queue with the given level and mask
	QueueVec(opts QueueVecOpts, level Level, mask Mask) QueueVec
}

// Group represents a logical grouping of metrics (e.g., "web", "database", "pipeline")
type Group interface {
	// Factory returns a factory for creating metrics in this group
	Factory() Factory

	// SetLevel sets the level for this group
	SetLevel(level Level)

	// SetMask sets the mask for this group
	SetMask(mask Mask)

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
	SetGlobalMask(mask Mask)

	// GlobalContext returns the global metrics context
	GlobalContext() Context
}
