package umami

import "time"

// Metric is the base interface for all metrics
type Metric interface {
	SetLevel(level Level)
	Name() string
	Help() string
	Type() MetricType
	Level() Level
}

type CompositeMetric interface {
	Metric
	Components() []Metric
}

type NoopMetric interface {
	Metric
	constructorOpts() any
}

// VecLabels is a type that represents a set partition keys to values
type VecLabels map[string]string

type BasicMetricOpts struct {
	FromComposite bool
}

type MetricInfo struct {
	Name string
	Help string
}

type CounterOpts struct {
	BasicMetricOpts
	MetricInfo
}

// Counter is a metric that counts occurrences. It only counts up.
type Counter interface {
	Metric

	// Inc increments the counter. Noop if disabled.
	Inc(ctx Context) error

	// Add adds the given value to the counter. Noop if disabled.
	Add(ctx Context, value float64) error
}

type CounterVecOpts struct {
	BasicMetricOpts
	MetricInfo
	Labels []string
}

// CounterVec is a metric that counts occurrences, partitioned by labels.
type CounterVec interface {
	Metric

	// Inc increments the counter for the given labels. Noop if disabled.
	Inc(ctx Context, labels VecLabels) error

	// Add adds the given value to the counter for the given labels. Noop if disabled.
	Add(ctx Context, value float64, labels VecLabels) error
}

type GaugeOpts struct {
	BasicMetricOpts
	MetricInfo
}

// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
type Gauge interface {
	Metric

	// Set sets the gauge to the given value. Noop if disabled.
	Set(ctx Context, value float64) error

	// Inc increments the gauge. Noop if disabled.
	Inc(ctx Context) error

	// Dec decrements the gauge. Noop if disabled.
	Dec(ctx Context) error

	// Add adds the given value to the gauge. Noop if disabled.
	Add(ctx Context, value float64) error
}

type GaugeVecOpts struct {
	BasicMetricOpts
	MetricInfo
	Labels []string
}

// GaugeVec is a metric that represents a collection of gauge values, partitioned by labels.
type GaugeVec interface {
	Metric

	// Set sets the gauge for the given labels to the given value. Noop if disabled.
	Set(ctx Context, value float64, labels VecLabels) error

	// Inc increments the gauge for the given labels. Noop if disabled.
	Inc(ctx Context, labels VecLabels) error

	// Dec decrements the gauge for the given labels. Noop if disabled.
	Dec(ctx Context, labels VecLabels) error

	// Add adds the given value to the gauge for the given labels. Noop if disabled.
	Add(ctx Context, value float64, labels VecLabels) error
}

type HistogramOpts struct {
	BasicMetricOpts
	MetricInfo
	Buckets []float64
}

// Histogram is a metric that represents a distribution of values.
type Histogram interface {
	Metric

	// Observe adds an observation to the histogram. Noop if disabled.
	Observe(ctx Context, value float64) error
}

type HistogramVecOpts struct {
	BasicMetricOpts
	MetricInfo
	Labels  []string
	Buckets []float64
}

// HistogramVec is a metric that represents a distribution of values, partitioned by labels.
type HistogramVec interface {
	Metric

	// Observe adds an observation to the histogram for the given labels. Noop if disabled.
	Observe(ctx Context, value float64, labels VecLabels) error
}

type SummaryOpts struct {
	BasicMetricOpts
	MetricInfo
	Objectives map[float64]float64
}

// Summary is a metric that provides quantiles of a distribution.
type Summary interface {
	Metric

	// Observe adds an observation to the summary. Noop if disabled.
	Observe(ctx Context, value float64) error

	// Quantile returns the value at the given quantile. Returns 0 if metric is disabled.
	Quantile(ctx Context, q float64) (float64, error)
}

type SummaryVecOpts struct {
	BasicMetricOpts
	MetricInfo
	Labels     []string
	Objectives map[float64]float64
}

// SummaryVec is a metric that provides quantiles of a distribution, partitioned by labels.
type SummaryVec interface {
	Metric

	// Observe adds an observation to the summary for the given labels. Noop if disabled.
	Observe(ctx Context, value float64, labels VecLabels) error

	// Quantile returns the value at the given quantile for the given labels. Returns 0 if metric is disabled.
	Quantile(ctx Context, q float64, labels VecLabels) (float64, error)
}

type TimerOpts struct {
	MetricInfo
	HistogramOpts HistogramOpts
}

// Timer is a metric that measures durations.
type Timer interface {
	CompositeMetric

	// Start returns a function that should be called when the operation completes
	// Returns a no-op function if metric is disabled
	Start(ctx Context) func()

	// Record records a duration. Noop if disabled.
	Record(ctx Context, duration time.Duration) error
}

type TimerVecOpts struct {
	MetricInfo
	HistogramVecOpts HistogramVecOpts
}

// TimerVec is a metric that measures durations, partitioned by labels.
type TimerVec interface {
	CompositeMetric

	// Start returns a function that should be called when the operation completes
	// Returns a no-op function if metric is disabled
	Start(ctx Context, labels VecLabels) func()

	// Record records a duration. Noop if disabled.
	Record(ctx Context, duration time.Duration, labels VecLabels) error
}

type CacheOpts struct {
	MetricInfo
	HitOpts  CounterOpts
	MissOpts CounterOpts
	SizeOpts GaugeOpts
}

// Cache is a metric that represents cache performance.
type Cache interface {
	CompositeMetric

	// Hit records a cache hit. Noop if disabled.
	Hit(ctx Context) error

	// Miss records a cache miss. Noop if disabled.
	Miss(ctx Context) error

	// SetSize sets the current cache size. Noop if disabled.
	SetSize(ctx Context, bytes int64) error
}

type CacheVecOpts struct {
	MetricInfo
	HitVecOpts  CounterVecOpts
	MissVecOpts CounterVecOpts
	SizeVecOpts GaugeVecOpts
}

// CacheVec is a metric that represents cache performance, partitioned by labels.
type CacheVec interface {
	CompositeMetric

	// Hit records a cache hit for the given labels. Noop if disabled.
	Hit(ctx Context, labels VecLabels) error

	// Miss records a cache miss for the given labels. Noop if disabled.
	Miss(ctx Context, labels VecLabels) error

	// SetSize sets the current cache size for the given labels. Noop if disabled.
	SetSize(ctx Context, bytes int64, labels VecLabels) error
}

type PoolOpts struct {
	MetricInfo
	ActiveOpts   GaugeOpts
	IdleOpts     GaugeOpts
	AcquiredOpts CounterOpts
	ReleasedOpts CounterOpts
}

// Pool is a metric that represents item pool utilization.
type Pool interface {
	CompositeMetric

	// SetActive sets the number of active items. Noop if disabled.
	SetActive(ctx Context, count int) error

	// SetIdle sets the number of idle items. Noop if disabled.
	SetIdle(ctx Context, count int) error

	// Acquired records an item acquisition. Noop if disabled.
	Acquired(ctx Context) error

	// Released records an item release. Noop if disabled.
	Released(ctx Context) error
}

type PoolVecOpts struct {
	MetricInfo
	ActiveVecOpts   GaugeVecOpts
	IdleVecOpts     GaugeVecOpts
	AcquiredVecOpts CounterVecOpts
	ReleasedVecOpts CounterVecOpts
}

// PoolVec is a metric that represents item pool utilization, partitioned by labels.
type PoolVec interface {
	CompositeMetric

	// SetActive sets the number of active items for the given labels. Noop if disabled.
	SetActive(ctx Context, count int, labels VecLabels) error

	// SetIdle sets the number of idle items for the given labels. Noop if disabled.
	SetIdle(ctx Context, count int, labels VecLabels) error

	// Acquired records an item acquisition for the given labels. Noop if disabled.
	Acquired(ctx Context, labels VecLabels) error

	// Released records an item release for the given labels. Noop if disabled.
	Released(ctx Context, labels VecLabels) error
}

type CircuitBreakerOpts struct {
	MetricInfo
	StateOpts   GaugeOpts
	SuccessOpts CounterOpts
	FailureOpts CounterOpts
}

// CircuitBreaker is a metric that represents the circuit breaker state
type CircuitBreaker interface {
	CompositeMetric

	// SetState sets the circuit breaker state. Noop if disabled.
	SetState(ctx Context, state CircuitBreakerState) error

	// Success records a successful operation. Noop if disabled.
	Success(ctx Context) error

	// Failure records a failed operation. Noop if disabled.
	Failure(ctx Context) error
}

type CircuitBreakerVecOpts struct {
	MetricInfo
	StateVecOpts   GaugeVecOpts
	SuccessVecOpts CounterVecOpts
	FailureVecOpts CounterVecOpts
}

// CircuitBreakerVec is a metric that represents the circuit breaker state, partitioned by labels.
type CircuitBreakerVec interface {
	CompositeMetric

	// SetState sets the circuit breaker state for the given labels. Noop if disabled.
	SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error

	// Success records a successful operation for the given labels. Noop if disabled.
	Success(ctx Context, labels VecLabels) error

	// Failure records a failed operation for the given labels. Noop if disabled.
	Failure(ctx Context, labels VecLabels) error
}

type QueueOpts struct {
	MetricInfo
	DepthOpts    GaugeOpts
	EnqueuedOpts CounterOpts
	DequeuedOpts CounterOpts
	WaitTimeOpts HistogramOpts
}

// Queue is a metric that represents queue statistics.
type Queue interface {
	CompositeMetric

	// SetDepth sets the current queue depth. Noop if disabled.
	SetDepth(ctx Context, depth int) error

	// Enqueued records an item being enqueued. Noop if disabled.
	Enqueued(ctx Context) error

	// Dequeued records an item being dequeued. Noop if disabled.
	Dequeued(ctx Context) error

	// SetWaitTime records how long items wait in the queue. Noop if disabled.
	SetWaitTime(ctx Context, duration time.Duration) error
}

type QueueVecOpts struct {
	MetricInfo
	DepthVecOpts    GaugeVecOpts
	EnqueuedVecOpts CounterVecOpts
	DequeuedVecOpts CounterVecOpts
	WaitTimeVecOpts HistogramVecOpts
}

// QueueVec is a metric that represents queue statistics, partitioned by labels.
type QueueVec interface {
	CompositeMetric

	// SetDepth sets the current queue depth for the given labels. Noop if disabled.
	SetDepth(ctx Context, depth int, labels VecLabels) error

	// Enqueued records an item being enqueued for the given labels. Noop if disabled.
	Enqueued(ctx Context, labels VecLabels) error

	// Dequeued records an item being dequeued for the given labels. Noop if disabled.
	Dequeued(ctx Context, labels VecLabels) error

	// SetWaitTime records how long items wait in the queue for the given labels. Noop if disabled.
	SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error
}
