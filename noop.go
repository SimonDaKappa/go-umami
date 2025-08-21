package umami

// No-op implementations for disabled metrics
// These provide zero-cost operations when metrics are disabled
//
// Generally, these are returned by factory functions when either
// the level or mask has disabled the requested metric.

import "time"

// noopCounter implements [Counter] interface with no-op operations
type noopCounter struct{}

func (n *noopCounter) Inc(ctx Context) error                { return nil }
func (n *noopCounter) Add(ctx Context, value float64) error { return nil }

// noopCounterVec implements [CounterVec] interface with no-op operations
type noopCounterVec struct{}

func (n *noopCounterVec) Inc(ctx Context, labels VecLabels) error                { return nil }
func (n *noopCounterVec) Add(ctx Context, value float64, labels VecLabels) error { return nil }

// noopGauge implements [Gauge] interface with no-op operations
type noopGauge struct{}

func (n *noopGauge) Set(ctx Context, value float64) error { return nil }
func (n *noopGauge) Inc(ctx Context) error                { return nil }
func (n *noopGauge) Dec(ctx Context) error                { return nil }
func (n *noopGauge) Add(ctx Context, value float64) error { return nil }

// noopGaugeVec implements [GaugeVec] interface with no-op operations
type noopGaugeVec struct{}

func (n *noopGaugeVec) Set(ctx Context, value float64, labels VecLabels) error { return nil }
func (n *noopGaugeVec) Add(ctx Context, value float64, labels VecLabels) error { return nil }
func (n *noopGaugeVec) Inc(ctx Context, labels VecLabels) error                { return nil }
func (n *noopGaugeVec) Dec(ctx Context, labels VecLabels) error                { return nil }

// noopHistogram implements [Histogram] interface with no-op operations
type noopHistogram struct{}

func (n *noopHistogram) Observe(ctx Context, value float64) error { return nil }
func (n *noopHistogram) Time(ctx Context, fn func() error) error  { return fn() }

// noopHistogramVec implements [HistogramVec] interface with no-op operations
type noopHistogramVec struct{}

func (n *noopHistogramVec) Observe(ctx Context, value float64, labels VecLabels) error { return nil }
func (n *noopHistogramVec) Time(ctx Context, fn func() error, labels VecLabels) error  { return fn() }

// noopSummary implements [Summary] interface with no-op operations
type noopSummary struct{}

func (n *noopSummary) Observe(ctx Context, value float64) error         { return nil }
func (n *noopSummary) Quantile(ctx Context, q float64) (float64, error) { return 0, nil }

// noopSummaryVec implements [SummaryVec] interface with no-op operations
type noopSummaryVec struct{}

func (n *noopSummaryVec) Observe(ctx Context, value float64, labels VecLabels) error { return nil }
func (n *noopSummaryVec) Quantile(ctx Context, q float64, labels VecLabels) (float64, error) {
	return 0, nil
}

// noopTimer implements [Timer] interface with no-op operations
type noopTimer struct{}

func (n *noopTimer) Start(ctx Context) func()                         { return func() {} }
func (n *noopTimer) Record(ctx Context, duration time.Duration) error { return nil }

// noopTimerVec implements [TimerVec] interface with no-op operations
type noopTimerVec struct{}

func (n *noopTimerVec) Start(ctx Context, labels VecLabels) func() { return func() {} }
func (n *noopTimerVec) Record(ctx Context, duration time.Duration, labels VecLabels) error {
	return nil
}

// noopCache implements [Cache] interface with no-op operations
type noopCache struct{}

func (n *noopCache) Hit(ctx Context) error                  { return nil }
func (n *noopCache) Miss(ctx Context) error                 { return nil }
func (n *noopCache) SetSize(ctx Context, bytes int64) error { return nil }

// noopCacheVec implements [CacheVec] interface with no-op operations
type noopCacheVec struct{}

func (n *noopCacheVec) Hit(ctx Context, labels VecLabels) error                  { return nil }
func (n *noopCacheVec) Miss(ctx Context, labels VecLabels) error                 { return nil }
func (n *noopCacheVec) SetSize(ctx Context, bytes int64, labels VecLabels) error { return nil }

// noopPool implements [Pool] interface with no-op operations
type noopPool struct{}

func (n *noopPool) SetActive(ctx Context, count int) error { return nil }
func (n *noopPool) SetIdle(ctx Context, count int) error   { return nil }
func (n *noopPool) Acquired(ctx Context) error             { return nil }
func (n *noopPool) Released(ctx Context) error             { return nil }

// noopPoolVec implements [PoolVec] interface with no-op operations
type noopPoolVec struct{}

func (n *noopPoolVec) SetActive(ctx Context, count int, labels VecLabels) error { return nil }
func (n *noopPoolVec) SetIdle(ctx Context, count int, labels VecLabels) error   { return nil }
func (n *noopPoolVec) Acquired(ctx Context, labels VecLabels) error             { return nil }
func (n *noopPoolVec) Released(ctx Context, labels VecLabels) error             { return nil }

// noopCircuitBreaker implements [CircuitBreaker] interface with no-op operations
type noopCircuitBreaker struct{}

func (n *noopCircuitBreaker) SetState(ctx Context, state CircuitBreakerState) error { return nil }
func (n *noopCircuitBreaker) Success(ctx Context) error                             { return nil }
func (n *noopCircuitBreaker) Failure(ctx Context) error                             { return nil }

// noopCircuitBreakerVec implements [CircuitBreakerVec] interface with no-op operations
type noopCircuitBreakerVec struct{}

func (n *noopCircuitBreakerVec) Success(ctx Context, labels VecLabels) error { return nil }
func (n *noopCircuitBreakerVec) Failure(ctx Context, labels VecLabels) error { return nil }
func (n *noopCircuitBreakerVec) SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error {
	return nil
}

// noopQueue implements [Queue] interface with no-op operations
type noopQueue struct{}

func (n *noopQueue) SetDepth(ctx Context, depth int) error                 { return nil }
func (n *noopQueue) Enqueued(ctx Context) error                            { return nil }
func (n *noopQueue) Dequeued(ctx Context) error                            { return nil }
func (n *noopQueue) SetWaitTime(ctx Context, duration time.Duration) error { return nil }

// noopQueueVec implements [QueueVec] interface with no-op operations
type noopQueueVec struct{}

func (n *noopQueueVec) Enqueued(ctx Context, labels VecLabels) error            { return nil }
func (n *noopQueueVec) Dequeued(ctx Context, labels VecLabels) error            { return nil }
func (n *noopQueueVec) SetDepth(ctx Context, depth int, labels VecLabels) error { return nil }
func (n *noopQueueVec) SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error {
	return nil
}

// Sanity checks for interfaces
var (
	_noopCounter           Counter           = (*noopCounter)(nil)
	_noopCounterVec        CounterVec        = (*noopCounterVec)(nil)
	_noopGauge             Gauge             = (*noopGauge)(nil)
	_noopGaugeVec          GaugeVec          = (*noopGaugeVec)(nil)
	_noopHistogram         Histogram         = (*noopHistogram)(nil)
	_noopHistogramVec      HistogramVec      = (*noopHistogramVec)(nil)
	_noopSummary           Summary           = (*noopSummary)(nil)
	_noopSummaryVec        SummaryVec        = (*noopSummaryVec)(nil)
	_noopTimer             Timer             = (*noopTimer)(nil)
	_noopTimerVec          TimerVec          = (*noopTimerVec)(nil)
	_noopCache             Cache             = (*noopCache)(nil)
	_noopCacheVec          CacheVec          = (*noopCacheVec)(nil)
	_noopPool              Pool              = (*noopPool)(nil)
	_noopPoolVec           PoolVec           = (*noopPoolVec)(nil)
	_noopCircuitBreaker    CircuitBreaker    = (*noopCircuitBreaker)(nil)
	_noopCircuitBreakerVec CircuitBreakerVec = (*noopCircuitBreakerVec)(nil)
	_noopQueue             Queue             = (*noopQueue)(nil)
	_noopQueueVec          QueueVec          = (*noopQueueVec)(nil)
)
