package umami

//--------------------------------------------------------------------------------
// File: base_metrics.go
//
// This file contains the base implementations of all metrics defined in the
// umami package. These implementations are designed to be composed with
// adapters provided by various backends to create concrete metric types.
//
// The base implementations handle level-based filtering and any common logic
// that is not specific to a particular backend.
//
// Almost all base metric impls will be wrapped in switchable wrappers to
// allow for dynamic run-time conversion between noop and real implementations
//
// Each base metric may be either a "basic" metric that directly composes a
// backend adapter, or a "composite" metric that composes one or more basic
// or other composite metrics to provide extended functionality.
//
// basic metrics include:
// - baseCounter
// - baseCounterVec
// - baseGauge
// - baseGaugeVec
// - baseHistogram
// - baseHistogramVec
// - baseSummary
// - baseSummaryVec
//
// Composite metrics include:
// - baseTimer (composes a Histogram)
// - baseTimerVec (composes a HistogramVec)
// - baseCache (composes Counters and a Gauge)
// - baseCacheVec (composes CounterVecs and a GaugeVec)
// - basePool (composes Gauges and Counters)
// - basePoolVec (composes GaugeVecs and CounterVecs)
// - baseCircuitBreaker (composes a Gauge and Counters)
// - baseCircuitBreakerVec (composes a GaugeVec and CounterVecs)
// - baseQueue (composes a Gauge, Counters, and a Histogram)
// - baseQueueVec (composes a GaugeVec, CounterVecs, and a HistogramVec)
//--------------------------------------------------------------------------------

import "time"

//--------------------------------------------------------------------------------
// Basic Base Metric Implementations
//
// Basic Base Metrics always compose an adapter, and store the level they
// were created at. Any call to a method on a [BasicMetric] will noop
// if the level in the passed context is below that of the instances.
//--------------------------------------------------------------------------------

// baseMetric provides common fields and methods for all metrics.
//
// It provides the default implementation for [Metric] interface methods.
// To use, embed this struct in your metric implementation and implement
// its specific methods.
//
// Note: Composite metrics must instead embed [baseCompositeMetric] or
// override the [baseMetric.SetLevel] method to propagate level changes,
// to composed basic metrics.
type baseMetric struct {
	level Level
	name  string
	help  string
}

func (b *baseMetric) Name() string {
	return b.name
}

func (b *baseMetric) Help() string {
	return b.help
}

func (b *baseMetric) SetLevel(level Level) {
	b.level = level
}

func (b *baseMetric) Type() MetricType {
	return MetricTypeBasic
}

func (b *baseMetric) Level() Level {
	return b.level
}

// baseCompositeMetric provides common fields and methods for composite metrics.
//
// It embeds [baseMetric] to inherit common functionality, but overrides
// the [SetLevel] method to propagate level changes to composed metrics.
//
// Note: Inheriting structs must implement the [CompositeMetric.Components] method to
// return the composed metrics for level propagation by [baseCompositeMetric.SetLevel].
type baseCompositeMetric struct {
	baseMetric
}

func (b *baseCompositeMetric) SetLevel(level Level) {
	for _, component := range b.Components() {
		component.SetLevel(level)
	}
}

func (b *baseCompositeMetric) Type() MetricType {
	return MetricTypeComposite
}

// Default implementation returns nil. MUST be overridden by inheriting structs.
func (b *baseCompositeMetric) Components() []Metric {
	return nil
}

//--------------------------------------------------------------------------------
// Basic Base Metric Implementations
//
// Required overrides:
// - NONE
//--------------------------------------------------------------------------------

type baseCounter struct {
	baseMetric
	adapter CounterAdapter
}

func (c *baseCounter) Inc(ctx Context) error {
	if !ctx.Enabled(c.level) {
		return nil
	}
	return c.adapter.Inc()
}

func (c *baseCounter) Add(ctx Context, value float64) error {
	if !ctx.Enabled(c.level) {
		return nil
	}
	return c.adapter.Add(value)
}

type baseCounterVec struct {
	baseMetric
	adapter CounterVecAdapter
}

func (cv *baseCounterVec) Inc(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(cv.level) {
		return nil
	}
	return cv.adapter.Inc(labels)
}

func (cv *baseCounterVec) Add(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(cv.level) {
		return nil
	}
	return cv.adapter.Add(value, labels)
}

type baseGauge struct {
	baseMetric
	adapter GaugeAdapter
}

func (g *baseGauge) Set(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) {
		return nil
	}
	return g.adapter.Set(value)
}

func (g *baseGauge) Inc(ctx Context) error {
	if !ctx.Enabled(g.level) {
		return nil
	}
	return g.adapter.Inc()
}

func (g *baseGauge) Dec(ctx Context) error {
	if !ctx.Enabled(g.level) {
		return nil
	}
	return g.adapter.Dec()
}

func (g *baseGauge) Add(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) {
		return nil
	}
	return g.adapter.Add(value)
}

type baseGaugeVec struct {
	baseMetric
	adapter GaugeVecAdapter
}

func (gv *baseGaugeVec) Set(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(gv.level) {
		return nil
	}
	return gv.adapter.Set(value, labels)
}

func (gv *baseGaugeVec) Inc(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(gv.level) {
		return nil
	}
	return gv.adapter.Inc(labels)
}

func (gv *baseGaugeVec) Dec(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(gv.level) {
		return nil
	}
	return gv.adapter.Dec(labels)
}

func (gv *baseGaugeVec) Add(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(gv.level) {
		return nil
	}
	return gv.adapter.Add(value, labels)
}

type baseHistogram struct {
	baseMetric
	adapter HistogramAdapter
}

func (h *baseHistogram) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(h.level) {
		return nil
	}
	return h.adapter.Observe(value)
}

// histogram wraps a HistogramBackend and implements early return
type baseHistogramVec struct {
	baseMetric
	adapter HistogramVecAdapter
}

func (hv *baseHistogramVec) Observe(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(hv.level) {
		return nil
	}
	return hv.adapter.Observe(value, labels)
}

type baseSummary struct {
	baseMetric
	adapter SummaryAdapter
}

func (s *baseSummary) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(s.level) {
		return nil
	}

	return s.adapter.Observe(value)
}

func (s *baseSummary) Quantile(ctx Context, q float64) (float64, error) {
	if !ctx.Enabled(s.level) {
		return 0, nil
	}

	return s.adapter.Quantile(q)
}

type baseSummaryVec struct {
	baseMetric
	adapter SummaryVecAdapater
}

func (sv *baseSummaryVec) Observe(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(sv.level) {
		return nil
	}
	return sv.adapter.Observe(value, labels)
}

func (sv *baseSummaryVec) Quantile(ctx Context, q float64, labels VecLabels) (float64, error) {
	if !ctx.Enabled(sv.level) {
		return 0, nil
	}
	return sv.adapter.Quantile(q, labels)
}

//--------------------------------------------------------------------------------
// Composite Base Metric Implementations.
//
// Composite Base Metrics compose basic (or other composite!) metrics to provide
// extended functionality and more complex behavior.
//
// Note: Composite Metrics must override the [baseMetric.SetLevel] method to
// propagate level changes
//
// Composite Base Metrics also store the level they were created at, and noop
// any method calls if the level in the passed context is below that of the
// instance. They should be designed in such a way that no composed metrics level
// supercedes that of the composite. That is, if the composite is disabled, all
// levels of composed metrics will be short-circuited and forcibly treated as
// disabled
//
// Required overrides:
// - Components() to return composed metrics for level propagation
//--------------------------------------------------------------------------------

type baseTimer struct {
	baseCompositeMetric
	histogram Histogram
}

func (t *baseTimer) Start(ctx Context) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		t.histogram.Observe(ctx, duration.Seconds())
	}
}

func (t *baseTimer) Record(ctx Context, duration time.Duration) error {
	return t.histogram.Observe(ctx, duration.Seconds())
}

func (t *baseTimer) Components() []Metric {
	return []Metric{t.histogram}
}

type baseTimerVec struct {
	baseCompositeMetric
	histogramVec HistogramVec
}

func (tv *baseTimerVec) Start(ctx Context, labels VecLabels) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		tv.histogramVec.Observe(ctx, duration.Seconds(), labels)
	}
}

func (tv *baseTimerVec) Record(ctx Context, duration time.Duration, labels VecLabels) error {
	return tv.histogramVec.Observe(ctx, duration.Seconds(), labels)
}

func (tv *baseTimerVec) Components() []Metric {
	return []Metric{tv.histogramVec}
}

type baseCache struct {
	baseCompositeMetric
	hits   Counter
	misses Counter
	size   Gauge
}

func (c *baseCache) Hit(ctx Context) error {
	return c.hits.Inc(ctx)
}

func (c *baseCache) Miss(ctx Context) error {
	return c.misses.Inc(ctx)
}

func (c *baseCache) SetSize(ctx Context, bytes int64) error {
	return c.size.Set(ctx, float64(bytes))
}

func (c *baseCache) Components() []Metric {
	return []Metric{c.hits, c.misses, c.size}
}

type baseCacheVec struct {
	baseCompositeMetric
	hits   CounterVec
	misses CounterVec
	size   GaugeVec
}

func (cv *baseCacheVec) Hit(ctx Context, labels VecLabels) error {
	return cv.hits.Inc(ctx, labels)
}

func (cv *baseCacheVec) Miss(ctx Context, labels VecLabels) error {
	return cv.misses.Inc(ctx, labels)
}

func (cv *baseCacheVec) SetSize(ctx Context, bytes int64, labels VecLabels) error {
	return cv.size.Set(ctx, float64(bytes), labels)
}

func (cv *baseCacheVec) Components() []Metric {
	return []Metric{cv.hits, cv.misses, cv.size}
}

type basePool struct {
	baseCompositeMetric
	active   Gauge
	idle     Gauge
	acquired Counter
	released Counter
}

func (p *basePool) SetActive(ctx Context, count int) error {
	return p.active.Set(ctx, float64(count))
}

func (p *basePool) SetIdle(ctx Context, count int) error {
	return p.idle.Set(ctx, float64(count))
}

func (p *basePool) Acquired(ctx Context) error {
	return p.acquired.Inc(ctx)
}

func (p *basePool) Released(ctx Context) error {
	return p.released.Inc(ctx)
}

func (p *basePool) Components() []Metric {
	return []Metric{p.active, p.idle, p.acquired, p.released}
}

type basePoolVec struct {
	baseCompositeMetric
	active   GaugeVec
	idle     GaugeVec
	acquired CounterVec
	released CounterVec
}

func (pv *basePoolVec) SetActive(ctx Context, count int, labels VecLabels) error {
	return pv.active.Set(ctx, float64(count), labels)
}

func (pv *basePoolVec) SetIdle(ctx Context, count int, labels VecLabels) error {
	return pv.idle.Set(ctx, float64(count), labels)
}

func (pv *basePoolVec) Acquired(ctx Context, labels VecLabels) error {
	return pv.acquired.Inc(ctx, labels)
}

func (pv *basePoolVec) Released(ctx Context, labels VecLabels) error {
	return pv.released.Inc(ctx, labels)
}

func (pv *basePoolVec) Components() []Metric {
	return []Metric{pv.active, pv.idle, pv.acquired, pv.released}
}

type baseCircuitBreaker struct {
	baseCompositeMetric
	state     Gauge
	successes Counter
	failures  Counter
}

func (cb *baseCircuitBreaker) SetState(ctx Context, state CircuitBreakerState) error {
	var value float64
	switch state {
	case CircuitBreakerStateClosed, CircuitBreakerStateOpen, CircuitBreakerStateHalfOpen:
		value = float64(state)
	default:
		value = -1
	}
	return cb.state.Set(ctx, value)
}

func (cb *baseCircuitBreaker) Success(ctx Context) error {
	return cb.successes.Inc(ctx)
}

func (cb *baseCircuitBreaker) Failure(ctx Context) error {
	return cb.failures.Inc(ctx)
}

func (cb *baseCircuitBreaker) Components() []Metric {
	return []Metric{cb.state, cb.successes, cb.failures}
}

type CircuitBreakerState uint8

const (
	CircuitBreakerStateClosed   CircuitBreakerState = 0
	CircuitBreakerStateOpen     CircuitBreakerState = 1
	CircuitBreakerStateHalfOpen CircuitBreakerState = 2
)

type baseCircuitBreakerVec struct {
	baseCompositeMetric
	state     GaugeVec
	successes CounterVec
	failures  CounterVec
}

func (cbv *baseCircuitBreakerVec) SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error {
	var value float64
	switch state {
	case CircuitBreakerStateClosed, CircuitBreakerStateOpen, CircuitBreakerStateHalfOpen:
		value = float64(state)
	default:
		value = -1
	}
	return cbv.state.Set(ctx, value, labels)
}

func (cbv *baseCircuitBreakerVec) Success(ctx Context, labels VecLabels) error {
	return cbv.successes.Inc(ctx, labels)
}

func (cbv *baseCircuitBreakerVec) Failure(ctx Context, labels VecLabels) error {
	return cbv.failures.Inc(ctx, labels)
}

func (cbv *baseCircuitBreakerVec) Components() []Metric {
	return []Metric{cbv.state, cbv.successes, cbv.failures}
}

type baseQueue struct {
	baseCompositeMetric
	depth    Gauge
	enqueued Counter
	dequeued Counter
	waitTime Histogram
}

func (q *baseQueue) SetDepth(ctx Context, depth int) error {
	return q.depth.Set(ctx, float64(depth))
}

func (q *baseQueue) Enqueued(ctx Context) error {
	return q.enqueued.Inc(ctx)
}

func (q *baseQueue) Dequeued(ctx Context) error {
	return q.dequeued.Inc(ctx)
}

func (q *baseQueue) SetWaitTime(ctx Context, duration time.Duration) error {
	return q.waitTime.Observe(ctx, duration.Seconds())
}

func (q *baseQueue) Components() []Metric {
	return []Metric{q.depth, q.enqueued, q.dequeued, q.waitTime}
}

type baseQueueVec struct {
	baseCompositeMetric
	depth    GaugeVec
	enqueued CounterVec
	dequeued CounterVec
	waitTime HistogramVec
}

func (qv *baseQueueVec) SetDepth(ctx Context, depth int, labels VecLabels) error {
	return qv.depth.Set(ctx, float64(depth), labels)
}

func (qv *baseQueueVec) Enqueued(ctx Context, labels VecLabels) error {
	return qv.enqueued.Inc(ctx, labels)
}

func (qv *baseQueueVec) Dequeued(ctx Context, labels VecLabels) error {
	return qv.dequeued.Inc(ctx, labels)
}

func (qv *baseQueueVec) SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error {
	return qv.waitTime.Observe(ctx, duration.Seconds(), labels)
}

func (qv *baseQueueVec) Components() []Metric {
	return []Metric{qv.depth, qv.enqueued, qv.dequeued, qv.waitTime}
}

var (
	// Common Interface compliance checks
	__ctc_baseMetric          Metric          = (*baseMetric)(nil)
	__ctc_baseCompositeMetric CompositeMetric = (*baseCompositeMetric)(nil)

	// Basic metrics compliance checks
	__ctc_baseCounter      Counter      = (*baseCounter)(nil)
	__ctc_baseCounterVec   CounterVec   = (*baseCounterVec)(nil)
	__ctc_baseGauge        Gauge        = (*baseGauge)(nil)
	__ctc_baseGaugeVec     GaugeVec     = (*baseGaugeVec)(nil)
	__ctc_baseHistogram    Histogram    = (*baseHistogram)(nil)
	__ctc_baseHistogramVec HistogramVec = (*baseHistogramVec)(nil)
	__ctc_baseSummary      Summary      = (*baseSummary)(nil)
	__ctc_baseSummaryVec   SummaryVec   = (*baseSummaryVec)(nil)

	// Composite metrics compliance checks
	__ctc_baseTimer             Timer             = (*baseTimer)(nil)
	__ctc_baseTimerVec          TimerVec          = (*baseTimerVec)(nil)
	__ctc_baseCache             Cache             = (*baseCache)(nil)
	__ctc_baseCacheVec          CacheVec          = (*baseCacheVec)(nil)
	__ctc_basePool              Pool              = (*basePool)(nil)
	__ctc_basePoolVec           PoolVec           = (*basePoolVec)(nil)
	__ctc_baseCircuitBreaker    CircuitBreaker    = (*baseCircuitBreaker)(nil)
	__ctc_baseCircuitBreakerVec CircuitBreakerVec = (*baseCircuitBreakerVec)(nil)
	__ctc_baseQueue             Queue             = (*baseQueue)(nil)
	__ctc_baseQueueVec          QueueVec          = (*baseQueueVec)(nil)
)
