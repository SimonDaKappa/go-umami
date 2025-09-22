package umami

//--------------------------------------------------------------------------------
// File: group.go
//
// This file contains the definition and implementation of the [Group] and [Factory]
// interfaces for the umami metrics library.
//--------------------------------------------------------------------------------

import (
	"sync"
)

//--------------------------------------------------------------------------------
// Interfaces
//--------------------------------------------------------------------------------

// Group represents a logical grouping of metrics (e.g., "web", "database", "pipeline")
//
// Groups have their own [Factory] for creating metrics, and their own [Context]
// for level checks. Each group tracks its [Metric]s and [CompositeMetric]s.
type Group interface {
	Factory

	// SetLevel sets the minimum level that this factory will create metrics for.
	// Any metrics requested below this level will be no-op implementations.
	//
	// This will update every metric created by this factory to the new level.
	//
	// Optionally, you can provide a flag to replace existing no-op metrics with
	// real implementations if they are now enabled by the new level.
	SetGroupLevel(level Level, opts LevelOpts)

	// Context returns a context for this group
	Context() Context

	Metric(name string) Metric
}

// Factory creates metrics with the appropriate [Level]
//
// A factory may create any type of metric implementation, but
// will typically create [switchableMetric]s that can switch between
// the corresponding base adapter wrapper and a noop implementation.
type Factory interface {
	// Counter creates a counter with the given level and mask
	Counter(opts CounterOpts, level Level) Counter

	// CounterVec creates a label-vectorized counter with the given level and mask
	CounterVec(opts CounterVecOpts, level Level) CounterVec

	// Gauge creates a gauge with the given level and mask
	Gauge(opts GaugeOpts, level Level) Gauge

	// GaugeVec creates a label-vectorized gauge with the given level and mask
	GaugeVec(opts GaugeVecOpts, level Level) GaugeVec

	// Histogram creates a histogram with the given level and mask
	Histogram(opts HistogramOpts, level Level) Histogram

	// HistogramVec creates a label-vectorized histogram with the given level and mask
	HistogramVec(opts HistogramVecOpts, level Level) HistogramVec

	// Summary creates a summary with the given level and mask
	Summary(opts SummaryOpts, level Level) Summary

	// SummaryVec creates a label-vectorized summary with the given level and mask
	SummaryVec(opts SummaryVecOpts, level Level) SummaryVec

	//--------------------------------------------------------------------------------
	// Composite Metrics
	//--------------------------------------------------------------------------------

	// Timer creates a timer with the given level and mask
	Timer(opts TimerOpts, level Level) Timer

	// TimerVec creates a label-vectorized timer with the given level and mask
	TimerVec(opts TimerVecOpts, level Level) TimerVec

	// Cache creates cache metrics with the given level and mask
	Cache(opts CacheOpts, level Level) Cache

	// CacheVec creates a label-vectorized cache with the given level and mask
	CacheVec(opts CacheVecOpts, level Level) CacheVec

	// Pool creates pool metrics with the given level and mask
	Pool(opts PoolOpts, level Level) Pool

	// PoolVec creates a label-vectorized pool with the given level and mask
	PoolVec(opts PoolVecOpts, level Level) PoolVec

	// CircuitBreaker creates circuit breaker metrics with the given level and mask
	CircuitBreaker(opts CircuitBreakerOpts, level Level) CircuitBreaker

	// CircuitBreakerVec creates a label-vectorized circuit breaker with the given level and mask
	CircuitBreakerVec(opts CircuitBreakerVecOpts, level Level) CircuitBreakerVec

	// Queue creates queue metrics with the given level and mask
	Queue(opts QueueOpts, level Level) Queue

	// QueueVec creates a label-vectorized queue with the given level and mask
	QueueVec(opts QueueVecOpts, level Level) QueueVec
}

//--------------------------------------------------------------------------------
// Group Implementation
//--------------------------------------------------------------------------------

// group implements the [Group] interface
type group struct {
	mu         sync.RWMutex
	name       string
	backend    Backend
	basics     map[string]SwitchableMetric
	composites map[string]SwitchableMetric
	noops      map[string]MetricType
	minLevel   Level
}

func newGroup(backend Backend, name string, level Level) *group {

	return &group{
		name:       name,
		minLevel:   level,
		backend:    backend,
		basics:     make(map[string]SwitchableMetric),
		composites: make(map[string]SwitchableMetric),
		noops:      make(map[string]MetricType),
	}
}

func (g *group) SetGroupLevel(level Level, opts LevelOpts) {
	g.minLevel = level

	if opts.ReplaceNoops && level.Enabled(g.minLevel) {
		g.convertNoops()
	} else {
		for _, metric := range g.composites {
			metric.SetLevel(level)
		}

		for _, metric := range g.basics {
			metric.SetLevel(level)
		}
	}
}

// Context returns a context representation of this group
func (g *group) Context() Context {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return NewContext(g.minLevel)
}

func (g *group) Metric(name string) Metric {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, metric := range g.basics {
		if metric.Name() == name {
			return metric
		}
	}

	for _, metric := range g.composites {
		if metric.Name() == name {
			return metric
		}
	}

	return nil
}

//--------------------------------------------------------------------------------
// Basic Metric Factory Functions
//
// Basic Metrics (aka [Metric]) are fundamental metric types that only have a
// single adapter instance, and compose no other [Metric]s into themselves.
// Basic Metrics may be a component of a [CompositeMetric], in which case, the
// parent [ComposeMetric] manages the level of the basic metric, and the
// basic metric is created with the [PrimeOpts]. FromComposite] option set to true, so it
// is not tracked individually by the group.
//
// Note: If a metric with the same name is already tracked in [group.primes], that
// metric is returned instead of creating a new one.
//--------------------------------------------------------------------------------

// Counter creates a counter with the given level
func (g *group) Counter(opts CounterOpts, level Level) Counter {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(Counter)
		}
	}

	var impl Counter
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		impl = newNoopCounter(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		impl = &baseCounter{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.Counter(opts),
		}
	}

	switchable := newSwitchableCounter(impl, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

func (g *group) CounterVec(opts CounterVecOpts, level Level) CounterVec {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(CounterVec)
		}
	}

	var counterVec CounterVec
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		counterVec = newNoopCounterVec(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		counterVec = &baseCounterVec{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.CounterVec(opts),
		}
	}

	switchable := newSwitchableCounterVec(counterVec, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// Gauge creates a gauge with the given level
func (g *group) Gauge(opts GaugeOpts, level Level) Gauge {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(Gauge)
		}
	}

	var gauge Gauge
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		gauge = newNoopGauge(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		gauge = &baseGauge{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.Gauge(opts),
		}
	}

	switchable := newSwitchableGauge(gauge, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// GaugeVec creates a gauge vector with the given level
func (g *group) GaugeVec(opts GaugeVecOpts, level Level) GaugeVec {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(GaugeVec)
		}
	}

	var gaugeVec GaugeVec
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		gaugeVec = newNoopGaugeVec(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		gaugeVec = &baseGaugeVec{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.GaugeVec(opts),
		}
	}

	switchable := newSwitchableGaugeVec(gaugeVec, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// Histogram creates a histogram with the given level
func (g *group) Histogram(opts HistogramOpts, level Level) Histogram {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(Histogram)
		}
	}

	var histogram Histogram
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		histogram = newNoopHistogram(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		histogram = &baseHistogram{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.Histogram(opts),
		}
	}

	switchable := newSwitchableHistogram(histogram, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// HistogramVec creates a histogram vector with the given level
func (g *group) HistogramVec(opts HistogramVecOpts, level Level) HistogramVec {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(HistogramVec)
		}
	}

	var histogramVec HistogramVec
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		histogramVec = newNoopHistogramVec(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		histogramVec = &baseHistogramVec{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.HistogramVec(opts),
		}
	}

	switchable := newSwitchableHistogramVec(histogramVec, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// Summary creates a summary with the given level
func (g *group) Summary(opts SummaryOpts, level Level) Summary {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(Summary)
		}
	}

	var summary Summary
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		summary = newNoopSummary(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		summary = &baseSummary{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.Summary(opts),
		}
	}

	switchable := newSwitchableSummary(summary, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

// SummaryVec creates a summary vector with the given level
func (g *group) SummaryVec(opts SummaryVecOpts, level Level) SummaryVec {
	opts.Name = g.name + "_" + opts.Name

	if !opts.FromComposite {
		if m := g.getBasic(opts.Name); m != nil {
			return m.(SummaryVec)
		}
	}

	var summaryVec SummaryVec
	var isTrackedNoop bool

	if !level.Enabled(g.minLevel) {
		summaryVec = newNoopSummaryVec(opts, level)
		isTrackedNoop = !opts.FromComposite
	} else {
		summaryVec = &baseSummaryVec{
			baseMetric: baseMetric{
				name:  opts.Name,
				help:  opts.Help,
				level: level,
			},
			adapter: g.backend.SummaryVec(opts),
		}
	}

	switchable := newSwitchableSummaryVec(summaryVec, opts)

	if !opts.FromComposite {
		g.track(switchable, isTrackedNoop)
	}

	return switchable
}

//--------------------------------------------------------------------------------
// Composite Metric Factory Functions
//
// Composite metrics are built on top of prime metrics, and so are tracked, but,
// their prime components are not tracked individually. So, any call to
// [group.SetGroupLevel] on a composite metric will not duplicate the calls to setting the
// prime components levels.
//--------------------------------------------------------------------------------

// Timer creates a timer with the given level
func (g *group) Timer(opts TimerOpts, level Level) Timer {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(Timer)
	}

	var timer Timer
	var isTrackedNoop bool
	opts.HistogramOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		timer = newNoopTimer(opts, level)
		isTrackedNoop = true
	} else {
		timer = &baseTimer{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			histogram: g.Histogram(opts.HistogramOpts, level),
		}
	}

	switchable := newSwitchableTimer(timer, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// TimerVec creates a timer vector with the given level
func (g *group) TimerVec(opts TimerVecOpts, level Level) TimerVec {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(TimerVec)
	}

	var timerVec TimerVec
	var isTrackedNoop bool

	opts.HistogramVecOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		timerVec = newNoopTimerVec(opts, level)
		isTrackedNoop = true
	} else {
		timerVec = &baseTimerVec{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			histogramVec: g.HistogramVec(opts.HistogramVecOpts, level),
		}
	}

	switchable := newSwitchableTimerVec(timerVec, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// Cache creates cache metrics with the given level
func (g *group) Cache(opts CacheOpts, level Level) Cache {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(Cache)
	}

	var cache Cache
	var isTrackedNoop bool

	opts.HitOpts.FromComposite = true
	opts.MissOpts.FromComposite = true
	opts.SizeOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		cache = newNoopCache(opts, level)
		isTrackedNoop = true
	} else {
		cache = &baseCache{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			hits:   g.Counter(opts.HitOpts, level),
			misses: g.Counter(opts.MissOpts, level),
			size:   g.Gauge(opts.SizeOpts, level),
		}
	}

	switchable := newSwitchableCache(cache, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// CacheVec creates a cache vector with the given level
func (g *group) CacheVec(opts CacheVecOpts, level Level) CacheVec {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(CacheVec)
	}

	var cacheVec CacheVec
	var isTrackedNoop bool

	opts.HitVecOpts.FromComposite = true
	opts.MissVecOpts.FromComposite = true
	opts.SizeVecOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		cacheVec = newNoopCacheVec(opts, level)
		isTrackedNoop = true
	} else {
		cacheVec = &baseCacheVec{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			hits:   g.CounterVec(opts.HitVecOpts, level),
			misses: g.CounterVec(opts.MissVecOpts, level),
			size:   g.GaugeVec(opts.SizeVecOpts, level),
		}
	}

	switchable := newSwitchableCacheVec(cacheVec, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// Pool creates pool metrics with the given level
func (g *group) Pool(opts PoolOpts, level Level) Pool {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(Pool)
	}

	var pool Pool
	var isTrackedNoop bool

	opts.ActiveOpts.FromComposite = true
	opts.IdleOpts.FromComposite = true
	opts.AcquiredOpts.FromComposite = true
	opts.ReleasedOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		pool = newNoopPool(opts, level)
		isTrackedNoop = true
	} else {
		pool = &basePool{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			active:   g.Gauge(opts.ActiveOpts, level),
			idle:     g.Gauge(opts.IdleOpts, level),
			acquired: g.Counter(opts.AcquiredOpts, level),
			released: g.Counter(opts.ReleasedOpts, level),
		}
	}

	switchable := newSwitchablePool(pool, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// PoolVec creates a pool vector with the given level
func (g *group) PoolVec(opts PoolVecOpts, level Level) PoolVec {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(PoolVec)
	}

	var poolVec PoolVec
	var isTrackedNoop bool

	opts.ActiveVecOpts.FromComposite = true
	opts.IdleVecOpts.FromComposite = true
	opts.AcquiredVecOpts.FromComposite = true
	opts.ReleasedVecOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		poolVec = newNoopPoolVec(opts, level)
		isTrackedNoop = true
	} else {

		poolVec = &basePoolVec{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			active:   g.GaugeVec(opts.ActiveVecOpts, level),
			idle:     g.GaugeVec(opts.IdleVecOpts, level),
			acquired: g.CounterVec(opts.AcquiredVecOpts, level),
			released: g.CounterVec(opts.ReleasedVecOpts, level),
		}
	}

	switchable := newSwitchablePoolVec(poolVec, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// CircuitBreaker creates circuit breaker metrics with the given level
func (g *group) CircuitBreaker(opts CircuitBreakerOpts, level Level) CircuitBreaker {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(CircuitBreaker)
	}

	var circuitBreaker CircuitBreaker
	var isTrackedNoop bool

	opts.StateOpts.FromComposite = true
	opts.SuccessOpts.FromComposite = true
	opts.FailureOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		circuitBreaker = newNoopCircuitBreaker(opts, level)
		isTrackedNoop = true
	} else {
		circuitBreaker = &baseCircuitBreaker{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			state:     g.Gauge(opts.StateOpts, level),
			successes: g.Counter(opts.SuccessOpts, level),
			failures:  g.Counter(opts.FailureOpts, level),
		}
	}

	switchable := newSwitchableCircuitBreaker(circuitBreaker, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// CircuitBreakerVec creates a circuit breaker vector with the given level
func (g *group) CircuitBreakerVec(opts CircuitBreakerVecOpts, level Level) CircuitBreakerVec {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(CircuitBreakerVec)
	}

	var circuitBreakerVec CircuitBreakerVec
	var isTrackedNoop bool

	opts.StateVecOpts.FromComposite = true
	opts.SuccessVecOpts.FromComposite = true
	opts.FailureVecOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		circuitBreakerVec = newNoopCircuitBreakerVec(opts, level)
		isTrackedNoop = true
	} else {

		circuitBreakerVec = &baseCircuitBreakerVec{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			state:     g.GaugeVec(opts.StateVecOpts, level),
			successes: g.CounterVec(opts.SuccessVecOpts, level),
			failures:  g.CounterVec(opts.FailureVecOpts, level),
		}
	}

	switchable := newSwitchableCircuitBreakerVec(circuitBreakerVec, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// Queue creates queue metrics with the given level
func (g *group) Queue(opts QueueOpts, level Level) Queue {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(Queue)
	}

	var queue Queue
	var isTrackedNoop bool

	opts.DepthOpts.FromComposite = true
	opts.EnqueuedOpts.FromComposite = true
	opts.DequeuedOpts.FromComposite = true
	opts.WaitTimeOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		queue = newNoopQueue(opts, level)
		isTrackedNoop = true
	} else {

		queue = &baseQueue{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			depth:    g.Gauge(opts.DepthOpts, level),
			enqueued: g.Counter(opts.EnqueuedOpts, level),
			dequeued: g.Counter(opts.DequeuedOpts, level),
			waitTime: g.Histogram(opts.WaitTimeOpts, level),
		}
	}

	switchable := newSwitchableQueue(queue, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

// QueueVec creates a queue vector with the given level
func (g *group) QueueVec(opts QueueVecOpts, level Level) QueueVec {
	if m := g.getComposite(opts.Name); m != nil {
		return m.(QueueVec)
	}

	var queueVec QueueVec
	var isTrackedNoop bool

	opts.DepthVecOpts.FromComposite = true
	opts.EnqueuedVecOpts.FromComposite = true
	opts.DequeuedVecOpts.FromComposite = true
	opts.WaitTimeVecOpts.FromComposite = true

	if !level.Enabled(g.minLevel) {
		queueVec = newNoopQueueVec(opts, level)
		isTrackedNoop = true
	} else {

		queueVec = &baseQueueVec{
			baseCompositeMetric: baseCompositeMetric{
				baseMetric: baseMetric{
					name:  opts.Name,
					help:  opts.Help,
					level: level,
				},
			},
			depth:    g.GaugeVec(opts.DepthVecOpts, level),
			enqueued: g.CounterVec(opts.EnqueuedVecOpts, level),
			dequeued: g.CounterVec(opts.DequeuedVecOpts, level),
			waitTime: g.HistogramVec(opts.WaitTimeVecOpts, level),
		}
	}

	switchable := newSwitchableQueueVec(queueVec, opts)
	switchable.SetLevel(level)

	g.track(switchable, isTrackedNoop)

	return switchable
}

//--------------------------------------------------------------------------------
// Metric Tracking Helpers
//--------------------------------------------------------------------------------

// getBasic retrieves a basic metric by name from
// the group's tracking map, or nil if not found.
//
// It is safe for concurrent use and read locks [group.mu]
func (g *group) getBasic(name string) Metric {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if m, exist := g.basics[name]; exist {
		return m
	}

	return nil
}

// getComposite retrieves a composite metric by name from
// the group's tracking map, or nil if not found.
//
// It is safe for concurrent use and read locks [group.mu]
func (g *group) getComposite(name string) Metric {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if m, exist := g.composites[name]; exist {
		return m
	}

	return nil
}

// track adds a metric (basic or composite) to the group's corresponding
// tracking map ([group.basics] or [group.composites]).
//
// WARNING: this does not check for duplicates, so the caller must ensure
// that the metric is not already tracked.
//
// Optionally, if isTrackedNoop is true, its name is added to the
// [group.noop] map for potential later replacement.
func (g *group) track(metric SwitchableMetric, isTrackedNoop bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch metric.Type() {
	case MetricTypeBasic:
		g.basics[metric.Name()] = metric
	case MetricTypeComposite:
		g.composites[metric.Name()] = metric
	}

	if isTrackedNoop {
		g.noops[metric.Name()] = metric.Type()
	}
}

//--------------------------------------------------------------------------------
// Noop Conversion and Group Management
//--------------------------------------------------------------------------------

func (g *group) convertNoopPrime(metric NoopMetric) Metric {

	switch metric.(type) {
	case *noopCounter:
		return g.Counter(metric.constructorOpts().(CounterOpts), metric.Level())
	case *noopCounterVec:
		return g.CounterVec(metric.constructorOpts().(CounterVecOpts), metric.Level())
	case *noopGauge:
		return g.Gauge(metric.constructorOpts().(GaugeOpts), metric.Level())
	case *noopGaugeVec:
		return g.GaugeVec(metric.constructorOpts().(GaugeVecOpts), metric.Level())
	case *noopHistogram:
		return g.Histogram(metric.constructorOpts().(HistogramOpts), metric.Level())
	case *noopHistogramVec:
		return g.HistogramVec(metric.constructorOpts().(HistogramVecOpts), metric.Level())
	case *noopSummary:
		return g.Summary(metric.constructorOpts().(SummaryOpts), metric.Level())
	case *noopSummaryVec:
		return g.SummaryVec(metric.constructorOpts().(SummaryVecOpts), metric.Level())
	default:
		panic("can't convert unknown basic NoopMetric type")
	}
}

func (g *group) convertNoopComposite(metric CompositeMetric) CompositeMetric {

	for _, component := range metric.Components() {
		if noop, ok := component.(NoopMetric); ok {

			// $$$SIMON can we guarantee that components always returns a slice of
			// pointers to the actual components, so that this assignment works?
			component = g.convertNoopPrime(noop)
		}
	}

	return metric
}

func (g *group) convertNoops() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, class := range g.noops {
		switch class {
		case 1:
			// do something
		case 2:
			// do something
		}
	}
}
