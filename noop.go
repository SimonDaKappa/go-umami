package umami

//--------------------------------------------------------------------------------
// File: noop.go
//
// This file contains no-op implementations of all metric interfaces.
// These are used when metrics are disabled via level settings.
//
// Typically, the factory for a group will return these no-op implementations
// when the requested metric's level is not enabled for the group.
//
// Each noop metric contains the constructor data that would be necessary to
// convert it to a real metric if the factory level is later changed to enable it.
//
// Noop metrics come in two flavors:
// 1. NoopPrimeMetric: Implements PrimeMetric and has constructorOpts() many,
//    that can be used to construct the backends actual implementation.
// 2. NoopCompositeMetric: Extends NoopPrimeMetric and also implements
//    CompositeMetric, which has a Components() method to return any underlying prime
//    metrics. This is used to construct the composite metric's actual
//    implementation.
//--------------------------------------------------------------------------------

// noopFunc is an empty function that is used by metric operations that return closures,
// but the level is disabled for the metric instance.
func noopFunc() {}

// noopCounter implements [Counter] interface with no-op operations
type noopCounter struct {
	baseMetric
	copts CounterOpts
}

func newNoopCounter(opts CounterOpts, level Level) *noopCounter {
	return &noopCounter{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopCounter) Inc(ctx Context) error {
	return nil
}

func (n *noopCounter) Add(ctx Context, value float64) error {
	return nil
}

func (n *noopCounter) constructorOpts() any {
	return n.copts
}

// noopCounterVec implements [CounterVec] interface with no-op operations
type noopCounterVec struct {
	baseMetric
	copts CounterVecOpts
}

func newNoopCounterVec(opts CounterVecOpts, level Level) *noopCounterVec {
	return &noopCounterVec{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopCounterVec) Inc(ctx Context, labels VecLabels) error {
	return nil
}

func (n *noopCounterVec) Add(ctx Context, value float64, labels VecLabels) error {
	return nil
}

func (n *noopCounterVec) constructorOpts() any {
	return n.copts
}

// noopGauge implements [Gauge] interface with no-op operations
type noopGauge struct {
	baseMetric
	copts GaugeOpts
}

func newNoopGauge(opts GaugeOpts, level Level) *noopGauge {
	return &noopGauge{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopGauge) Set(ctx Context, value float64) error {
	return nil
}

func (n *noopGauge) Inc(ctx Context) error {
	return nil
}

func (n *noopGauge) Dec(ctx Context) error {
	return nil
}

func (n *noopGauge) Add(ctx Context, value float64) error {
	return nil
}

func (n *noopGauge) constructorOpts() any {
	return n.copts
}

// noopGaugeVec implements [GaugeVec] interface with no-op operations
type noopGaugeVec struct {
	baseMetric
	copts GaugeVecOpts
}

func newNoopGaugeVec(opts GaugeVecOpts, level Level) *noopGaugeVec {
	return &noopGaugeVec{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopGaugeVec) Set(ctx Context, value float64, labels VecLabels) error {
	return nil
}

func (n *noopGaugeVec) Add(ctx Context, value float64, labels VecLabels) error {
	return nil
}

func (n *noopGaugeVec) Inc(ctx Context, labels VecLabels) error {
	return nil
}

func (n *noopGaugeVec) Dec(ctx Context, labels VecLabels) error {
	return nil
}

func (n *noopGaugeVec) constructorOpts() any {
	return n.copts
}

// noopHistogram implements [Histogram] interface with no-op operations
type noopHistogram struct {
	baseMetric
	copts HistogramOpts
}

func newNoopHistogram(opts HistogramOpts, level Level) *noopHistogram {
	return &noopHistogram{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopHistogram) Observe(ctx Context, value float64) error {
	return nil
}

func (n *noopHistogram) constructorOpts() any {
	return n.copts
}

// noopHistogramVec implements [HistogramVec] interface with no-op operations
type noopHistogramVec struct {
	baseMetric
	copts HistogramVecOpts
}

func newNoopHistogramVec(opts HistogramVecOpts, level Level) *noopHistogramVec {
	return &noopHistogramVec{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopHistogramVec) Observe(ctx Context, value float64, labels VecLabels) error {
	return nil
}

func (n *noopHistogramVec) constructorOpts() any {
	return n.copts
}

// noopSummary implements [Summary] interface with no-op operations
type noopSummary struct {
	baseMetric
	copts SummaryOpts
}

func newNoopSummary(opts SummaryOpts, level Level) *noopSummary {
	return &noopSummary{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopSummary) Observe(ctx Context, value float64) error {
	return nil
}

func (n *noopSummary) Quantile(ctx Context, q float64) (float64, error) {
	return 0, nil
}

func (n *noopSummary) constructorOpts() any {
	return n.copts
}

// noopSummaryVec implements [SummaryVec] interface with no-op operations
type noopSummaryVec struct {
	baseMetric
	copts SummaryVecOpts
}

func newNoopSummaryVec(opts SummaryVecOpts, level Level) *noopSummaryVec {
	return &noopSummaryVec{
		baseMetric: baseMetric{
			name:  opts.Name,
			help:  opts.Help,
			level: level,
		},
		copts: opts,
	}
}

func (n *noopSummaryVec) Observe(ctx Context, value float64, labels VecLabels) error {
	return nil
}

func (n *noopSummaryVec) Quantile(ctx Context, q float64, labels VecLabels) (float64, error) {
	return 0, nil
}

func (n *noopSummaryVec) constructorOpts() any {
	return n.copts
}

// // noopTimer implements [Timer] interface with no-op operations
// type noopTimer struct {
// 	baseMetric
// 	opts      TimerOpts
// 	histogram *noopHistogram
// }

func newNoopTimer(opts TimerOpts, level Level) Timer {
	opts.HistogramOpts.FromComposite = true
	opts.HistogramOpts.Name = opts.Name + "_histogram"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseTimer{
		baseCompositeMetric: baseCompositeMetric{base},
		histogram:           newNoopHistogram(opts.HistogramOpts, level),
	}
}

// func (n *noopTimer) SetLevel(level Level) {
// 	n.level = level
// 	n.histogram.SetLevel(level)
// }

// func (n *noopTimer) Start(ctx Context) func() {
// 	return noopFunc
// }

// func (n *noopTimer) Record(ctx Context, duration time.Duration) error {
// 	return nil
// }

// func (n *noopTimer) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopTimer) Components() []Metric {
// 	return []Metric{
// 		n.histogram,
// 	}
// }

// // noopTimerVec implements [TimerVec] interface with no-op operations
// type noopTimerVec struct {
// 	baseMetric
// 	opts         TimerVecOpts
// 	histogramVec *noopHistogramVec
// }

func newNoopTimerVec(opts TimerVecOpts, level Level) TimerVec {
	opts.HistogramVecOpts.FromComposite = true
	opts.HistogramVecOpts.Name = opts.Name + "_histogram"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseTimerVec{
		baseCompositeMetric: baseCompositeMetric{base},
		histogramVec:        newNoopHistogramVec(opts.HistogramVecOpts, level),
	}
}

// func (n *noopTimerVec) SetLevel(level Level) {
// 	n.level = level
// 	n.histogramVec.SetLevel(level)
// }

// func (n *noopTimerVec) Start(ctx Context, labels VecLabels) func() {
// 	return noopFunc
// }

// func (n *noopTimerVec) Record(ctx Context, duration time.Duration, labels VecLabels) error {
// 	return nil
// }

// func (n *noopTimerVec) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopTimerVec) Components() []Metric {
// 	return []Metric{
// 		n.histogramVec,
// 	}
// }

// // noopCache implements [Cache] interface with no-op operations
// type noopCache struct {
// 	baseMetric
// 	opts   CacheOpts
// 	hits   *noopCounter
// 	misses *noopCounter
// 	size   *noopGauge
// }

func newNoopCache(opts CacheOpts, level Level) Cache {
	opts.HitOpts.FromComposite = true
	opts.HitOpts.Name = opts.Name + "_hit"
	opts.MissOpts.FromComposite = true
	opts.MissOpts.Name = opts.Name + "_miss"
	opts.SizeOpts.FromComposite = true
	opts.SizeOpts.Name = opts.Name + "_size"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseCache{
		baseCompositeMetric: baseCompositeMetric{base},
		hits:                newNoopCounter(opts.HitOpts, level),
		misses:              newNoopCounter(opts.MissOpts, level),
		size:                newNoopGauge(opts.SizeOpts, level),
	}
}

// func (n *noopCache) SetLevel(level Level) {
// 	n.level = level
// 	n.hits.SetLevel(level)
// 	n.misses.SetLevel(level)
// 	n.size.SetLevel(level)
// }

// func (n *noopCache) Hit(ctx Context) error {
// 	return nil
// }

// func (n *noopCache) Miss(ctx Context) error {
// 	return nil
// }

// func (n *noopCache) SetSize(ctx Context, bytes int64) error {
// 	return nil
// }

// func (n *noopCache) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopCache) Components() []Metric {
// 	return []Metric{
// 		n.hits, n.misses, n.size,
// 	}
// }

// // noopCacheVec implements [CacheVec] interface with no-op operations
// type noopCacheVec struct {
// 	baseMetric
// 	opts   CacheVecOpts
// 	hits   *noopCounterVec
// 	misses *noopCounterVec
// 	size   *noopGaugeVec
// }

func newNoopCacheVec(opts CacheVecOpts, level Level) CacheVec {
	opts.HitVecOpts.FromComposite = true
	opts.HitVecOpts.Name = opts.Name + "_hit"
	opts.MissVecOpts.FromComposite = true
	opts.MissVecOpts.Name = opts.Name + "_miss"
	opts.SizeVecOpts.FromComposite = true
	opts.SizeVecOpts.Name = opts.Name + "_size"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseCacheVec{
		baseCompositeMetric: baseCompositeMetric{base},
		hits:                newNoopCounterVec(opts.HitVecOpts, level),
		misses:              newNoopCounterVec(opts.MissVecOpts, level),
		size:                newNoopGaugeVec(opts.SizeVecOpts, level),
	}
}

// func (n *noopCacheVec) SetLevel(level Level) {
// 	n.level = level
// 	n.hits.SetLevel(level)
// 	n.misses.SetLevel(level)
// 	n.size.SetLevel(level)
// }

// func (n *noopCacheVec) Hit(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCacheVec) Miss(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCacheVec) SetSize(ctx Context, bytes int64, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCacheVec) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopCacheVec) Components() []Metric {
// 	return []Metric{
// 		n.hits, n.misses, n.size,
// 	}
// }

// // noopPool implements [Pool] interface with no-op operations
// type noopPool struct {
// 	name     string
// 	opts     PoolOpts
// 	level    Level
// 	active   *noopGauge
// 	idle     *noopGauge
// 	acquired *noopCounter
// 	released *noopCounter
// }

func newNoopPool(opts PoolOpts, level Level) Pool {
	opts.ActiveOpts.FromComposite = true
	opts.ActiveOpts.Name = opts.Name + "_active"
	opts.IdleOpts.FromComposite = true
	opts.IdleOpts.Name = opts.Name + "_idle"
	opts.AcquiredOpts.FromComposite = true
	opts.AcquiredOpts.Name = opts.Name + "_acquired"
	opts.ReleasedOpts.FromComposite = true
	opts.ReleasedOpts.Name = opts.Name + "_released"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &basePool{
		baseCompositeMetric: baseCompositeMetric{base},
		active:              newNoopGauge(opts.ActiveOpts, level),
		idle:                newNoopGauge(opts.IdleOpts, level),
		acquired:            newNoopCounter(opts.AcquiredOpts, level),
		released:            newNoopCounter(opts.ReleasedOpts, level),
	}
}

// func (n *noopPool) SetLevel(level Level) {
// 	n.level = level
// 	n.active.SetLevel(level)
// 	n.idle.SetLevel(level)
// 	n.acquired.SetLevel(level)
// 	n.released.SetLevel(level)
// }
// func (n *noopPool) SetActive(ctx Context, count int) error {
// 	return nil
// }

// func (n *noopPool) SetIdle(ctx Context, count int) error {
// 	return nil
// }

// func (n *noopPool) Acquired(ctx Context) error {
// 	return nil
// }

// func (n *noopPool) Released(ctx Context) error {
// 	return nil
// }

// func (n *noopPool) Name() string {
// 	return n.name
// }

// func (n *noopPool) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopPool) Components() []Metric {
// 	return []Metric{
// 		n.active, n.idle, n.acquired, n.released,
// 	}
// }

// // noopPoolVec implements [PoolVec] interface with no-op operations
// type noopPoolVec struct {
// 	name     string
// 	opts     PoolVecOpts
// 	level    Level
// 	active   *noopGaugeVec
// 	idle     *noopGaugeVec
// 	acquired *noopCounterVec
// 	released *noopCounterVec
// }

func newNoopPoolVec(opts PoolVecOpts, level Level) PoolVec {
	opts.ActiveVecOpts.FromComposite = true
	opts.ActiveVecOpts.Name = opts.Name + "_active"
	opts.IdleVecOpts.FromComposite = true
	opts.IdleVecOpts.Name = opts.Name + "_idle"
	opts.AcquiredVecOpts.FromComposite = true
	opts.AcquiredVecOpts.Name = opts.Name + "_acquired"
	opts.ReleasedVecOpts.FromComposite = true
	opts.ReleasedVecOpts.Name = opts.Name + "_released"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &basePoolVec{
		baseCompositeMetric: baseCompositeMetric{base},
		active:              newNoopGaugeVec(opts.ActiveVecOpts, level),
		idle:                newNoopGaugeVec(opts.IdleVecOpts, level),
		acquired:            newNoopCounterVec(opts.AcquiredVecOpts, level),
		released:            newNoopCounterVec(opts.ReleasedVecOpts, level),
	}
}

// func (n *noopPoolVec) SetLevel(level Level) {
// 	n.level = level
// 	n.active.SetLevel(level)
// 	n.idle.SetLevel(level)
// 	n.acquired.SetLevel(level)
// 	n.released.SetLevel(level)
// }

// func (n *noopPoolVec) SetActive(ctx Context, count int, labels VecLabels) error {
// 	return nil
// }

// func (n *noopPoolVec) SetIdle(ctx Context, count int, labels VecLabels) error {
// 	return nil
// }

// func (n *noopPoolVec) Acquired(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopPoolVec) Released(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopPoolVec) Name() string {
// 	return n.name
// // }

// func (n *noopPoolVec) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopPoolVec) Components() []Metric {
// 	return []Metric{
// 		n.active, n.idle, n.acquired, n.released,
// 	}
// }

// // noopCircuitBreaker implements [CircuitBreaker] interface with no-op operations
// type noopCircuitBreaker struct {
// 	name      string
// 	opts      CircuitBreakerOpts
// 	level     Level
// 	state     *noopGauge
// 	successes *noopCounter
// 	failures  *noopCounter
// }

func newNoopCircuitBreaker(opts CircuitBreakerOpts, level Level) CircuitBreaker {
	opts.StateOpts.FromComposite = true
	opts.StateOpts.Name = opts.Name + "_state"
	opts.SuccessOpts.FromComposite = true
	opts.SuccessOpts.Name = opts.Name + "_success"
	opts.FailureOpts.FromComposite = true
	opts.FailureOpts.Name = opts.Name + "_failure"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseCircuitBreaker{
		baseCompositeMetric: baseCompositeMetric{base},
		state:               newNoopGauge(opts.StateOpts, level),
		successes:           newNoopCounter(opts.SuccessOpts, level),
		failures:            newNoopCounter(opts.FailureOpts, level),
	}
}

// func (n *noopCircuitBreaker) SetLevel(level Level) {
// 	n.level = level
// 	n.state.SetLevel(level)
// 	n.successes.SetLevel(level)
// 	n.failures.SetLevel(level)
// }

// func (n *noopCircuitBreaker) SetState(ctx Context, state CircuitBreakerState) error {
// 	return nil
// }

// func (n *noopCircuitBreaker) Success(ctx Context) error {
// 	return nil
// }

// func (n *noopCircuitBreaker) Failure(ctx Context) error {
// 	return nil
// }

// func (n *noopCircuitBreaker) Name() string {
// 	return n.name
// }

// func (n *noopCircuitBreaker) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopCircuitBreaker) Components() []Metric {
// 	return []Metric{
// 		n.state, n.successes, n.failures,
// 	}
// }

// // noopCircuitBreakerVec implements [CircuitBreakerVec] interface with no-op operations
// type noopCircuitBreakerVec struct {
// 	name      string
// 	opts      CircuitBreakerVecOpts
// 	level     Level
// 	state     *noopGaugeVec
// 	successes *noopCounterVec
// 	failures  *noopCounterVec
// }

func newNoopCircuitBreakerVec(opts CircuitBreakerVecOpts, level Level) CircuitBreakerVec {
	opts.StateVecOpts.FromComposite = true
	opts.StateVecOpts.Name = opts.Name + "_state"
	opts.SuccessVecOpts.FromComposite = true
	opts.SuccessVecOpts.Name = opts.Name + "_success"
	opts.FailureVecOpts.FromComposite = true
	opts.FailureVecOpts.Name = opts.Name + "_failure"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseCircuitBreakerVec{
		baseCompositeMetric: baseCompositeMetric{base},
		state:               newNoopGaugeVec(opts.StateVecOpts, level),
		successes:           newNoopCounterVec(opts.SuccessVecOpts, level),
		failures:            newNoopCounterVec(opts.FailureVecOpts, level),
	}
}

// func (n *noopCircuitBreakerVec) SetLevel(level Level) {
// 	n.level = level
// 	n.state.SetLevel(level)
// 	n.successes.SetLevel(level)
// 	n.failures.SetLevel(level)
// }

// func (n *noopCircuitBreakerVec) Success(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCircuitBreakerVec) Failure(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCircuitBreakerVec) SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error {
// 	return nil
// }

// func (n *noopCircuitBreakerVec) Name() string {
// 	return n.name
// }

// func (n *noopCircuitBreakerVec) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopCircuitBreakerVec) Components() []Metric {
// 	return []Metric{
// 		n.state, n.successes, n.failures,
// 	}
// }

// // noopQueue implements [Queue] interface with no-op operations
// type noopQueue struct {
// 	name     string
// 	opts     QueueOpts
// 	level    Level
// 	depth    *noopGauge
// 	enqueued *noopCounter
// 	dequeued *noopCounter
// 	waitTime *noopHistogram
// }

func newNoopQueue(opts QueueOpts, level Level) Queue {
	opts.DepthOpts.FromComposite = true
	opts.DepthOpts.Name = opts.Name + "_depth"
	opts.EnqueuedOpts.FromComposite = true
	opts.EnqueuedOpts.Name = opts.Name + "_enqueued"
	opts.DequeuedOpts.FromComposite = true
	opts.DequeuedOpts.Name = opts.Name + "_dequeued"
	opts.WaitTimeOpts.FromComposite = true
	opts.WaitTimeOpts.Name = opts.Name + "_wait_time"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseQueue{
		baseCompositeMetric: baseCompositeMetric{base},
		depth:               newNoopGauge(opts.DepthOpts, level),
		enqueued:            newNoopCounter(opts.EnqueuedOpts, level),
		dequeued:            newNoopCounter(opts.DequeuedOpts, level),
		waitTime:            newNoopHistogram(opts.WaitTimeOpts, level),
	}
}

// func (n *noopQueue) SetLevel(level Level) {
// 	n.level = level
// 	n.depth.SetLevel(level)
// 	n.enqueued.SetLevel(level)
// 	n.dequeued.SetLevel(level)
// 	n.waitTime.SetLevel(level)
// }

// func (n *noopQueue) SetDepth(ctx Context, depth int) error {
// 	return nil
// }

// func (n *noopQueue) Enqueued(ctx Context) error {
// 	return nil
// }

// func (n *noopQueue) Dequeued(ctx Context) error {
// 	return nil
// }

// func (n *noopQueue) SetWaitTime(ctx Context, duration time.Duration) error {
// 	return nil
// }

// func (n *noopQueue) Name() string {
// 	return n.name
// }

// func (n *noopQueue) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopQueue) Components() []Metric {
// 	return []Metric{
// 		n.depth, n.enqueued, n.dequeued, n.waitTime,
// 	}
// }

// // noopQueueVec implements [QueueVec] interface with no-op operations
// type noopQueueVec struct {
// 	name     string
// 	opts     QueueVecOpts
// 	level    Level
// 	depth    *noopGaugeVec
// 	enqueued *noopCounterVec
// 	dequeued *noopCounterVec
// 	waitTime *noopHistogramVec
// }

func newNoopQueueVec(opts QueueVecOpts, level Level) QueueVec {
	opts.DepthVecOpts.FromComposite = true
	opts.DepthVecOpts.Name = opts.Name + "_depth"
	opts.EnqueuedVecOpts.FromComposite = true
	opts.EnqueuedVecOpts.Name = opts.Name + "_enqueued"
	opts.DequeuedVecOpts.FromComposite = true
	opts.DequeuedVecOpts.Name = opts.Name + "_dequeued"
	opts.WaitTimeVecOpts.FromComposite = true
	opts.WaitTimeVecOpts.Name = opts.Name + "_wait_time"

	base := baseMetric{
		name:  opts.Name,
		help:  opts.Help,
		level: level,
	}

	return &baseQueueVec{
		baseCompositeMetric: baseCompositeMetric{base},
		depth:               newNoopGaugeVec(opts.DepthVecOpts, level),
		enqueued:            newNoopCounterVec(opts.EnqueuedVecOpts, level),
		dequeued:            newNoopCounterVec(opts.DequeuedVecOpts, level),
		waitTime:            newNoopHistogramVec(opts.WaitTimeVecOpts, level),
	}
}

// func (n *noopQueueVec) SetLevel(level Level) {
// 	n.level = level
// 	n.depth.SetLevel(level)
// 	n.enqueued.SetLevel(level)
// 	n.dequeued.SetLevel(level)
// 	n.waitTime.SetLevel(level)
// }

// func (n *noopQueueVec) Enqueued(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopQueueVec) Dequeued(ctx Context, labels VecLabels) error {
// 	return nil
// }

// func (n *noopQueueVec) SetDepth(ctx Context, depth int, labels VecLabels) error {
// 	return nil
// }

// func (n *noopQueueVec) SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error {
// 	return nil
// }

// func (n *noopQueueVec) Name() string {
// 	return n.name
// }

// func (n *noopQueueVec) constructorOpts() any {
// 	return n.opts
// }

// func (n *noopQueueVec) Components() []Metric {
// 	return []Metric{
// 		n.depth, n.enqueued, n.dequeued, n.waitTime,
// 	}
// }

// Sanity checks for interfaces
var (
	// Metric interface checks
	__ctc_noopCounterIntf           Counter           = (*noopCounter)(nil)
	__ctc_noopCounterVecIntf        CounterVec        = (*noopCounterVec)(nil)
	__ctc_noopGaugeIntf             Gauge             = (*noopGauge)(nil)
	__ctc_noopGaugeVecIntf          GaugeVec          = (*noopGaugeVec)(nil)
	__ctc_noopHistogramIntf         Histogram         = (*noopHistogram)(nil)
	__ctc_noopHistogramVecIntf      HistogramVec      = (*noopHistogramVec)(nil)
	__ctc_noopSummaryIntf           Summary           = (*noopSummary)(nil)
	__ctc_noopSummaryVecIntf        SummaryVec        = (*noopSummaryVec)(nil)
	__ctc_noopTimerIntf             Timer             = newNoopTimer(TimerOpts{}, LevelDisabled)
	__ctc_noopTimerVecIntf          TimerVec          = newNoopTimerVec(TimerVecOpts{}, LevelDisabled)
	__ctc_noopCacheIntf             Cache             = newNoopCache(CacheOpts{}, LevelDisabled)
	__ctc_noopCacheVecIntf          CacheVec          = newNoopCacheVec(CacheVecOpts{}, LevelDisabled)
	__ctc_noopPoolIntf              Pool              = newNoopPool(PoolOpts{}, LevelDisabled)
	__ctc_noopPoolVecIntf           PoolVec           = newNoopPoolVec(PoolVecOpts{}, LevelDisabled)
	__ctc_noopCircuitBreakerIntf    CircuitBreaker    = newNoopCircuitBreaker(CircuitBreakerOpts{}, LevelDisabled)
	__ctc_noopCircuitBreakerVecIntf CircuitBreakerVec = newNoopCircuitBreakerVec(CircuitBreakerVecOpts{}, LevelDisabled)
	__ctc_noopQueueIntf             Queue             = newNoopQueue(QueueOpts{}, LevelDisabled)
	__ctc_noopQueueVecIntf          QueueVec          = newNoopQueueVec(QueueVecOpts{}, LevelDisabled)

	// Basic NoopMetric interface checks
	__ctc_noopCounterNoopBasic      NoopMetric = (*noopCounter)(nil)
	__ctc_noopCounterVecNoopBasic   NoopMetric = (*noopCounterVec)(nil)
	__ctc_noopGaugeNoopBasic        NoopMetric = (*noopGauge)(nil)
	__ctc_noopGaugeVecNoopBasic     NoopMetric = (*noopGaugeVec)(nil)
	__ctc_noopHistogramNoopBasic    NoopMetric = (*noopHistogram)(nil)
	__ctc_noopHistogramVecNoopBasic NoopMetric = (*noopHistogramVec)(nil)
	__ctc_noopSummaryNoopBasic      NoopMetric = (*noopSummary)(nil)
	__ctc_noopSummaryVecNoopBasic   NoopMetric = (*noopSummaryVec)(nil)

	// Composite interface checks
	__ctc_noopTimerNoopComposite             CompositeMetric = newNoopTimer(TimerOpts{}, LevelDisabled)
	__ctc_noopTimerVecNoopComposite          CompositeMetric = newNoopTimerVec(TimerVecOpts{}, LevelDisabled)
	__ctc_noopCacheNoopComposite             CompositeMetric = newNoopCache(CacheOpts{}, LevelDisabled)
	__ctc_noopCacheVecNoopComposite          CompositeMetric = newNoopCacheVec(CacheVecOpts{}, LevelDisabled)
	__ctc_noopPoolNoopComposite              CompositeMetric = newNoopPool(PoolOpts{}, LevelDisabled)
	__ctc_noopPoolVecNoopComposite           CompositeMetric = newNoopPoolVec(PoolVecOpts{}, LevelDisabled)
	__ctc_noopCircuitBreakerNoopComposite    CompositeMetric = newNoopCircuitBreaker(CircuitBreakerOpts{}, LevelDisabled)
	__ctc_noopCircuitBreakerVecNoopComposite CompositeMetric = newNoopCircuitBreakerVec(CircuitBreakerVecOpts{}, LevelDisabled)
	__ctc_noopQueueNoopComposite             CompositeMetric = newNoopQueue(QueueOpts{}, LevelDisabled)
	__ctc_noopQueueVecNoopComposite          CompositeMetric = newNoopQueueVec(QueueVecOpts{}, LevelDisabled)
)
