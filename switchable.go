package umami

//--------------------------------------------------------------------------------
// File: switchable_metrics.go
//
// This file contains switchable metric wrappers that solve the noop-to-real
// conversion problem. Instead of returning noop metrics directly, the factory
// returns these wrappers that can switch their internal implementation
// without changing their identity.
//
// When a level change occurs, the group can call switchImpl() on
// these wrappers to replace the internal metric while preserving user references.
//--------------------------------------------------------------------------------

import (
	"sync"
	"time"
)

//--------------------------------------------------------------------------------
// Switchable Interface
//
// All switchable metrics implement this interface to allow the group to
// update their internal implementation when levels change.
//--------------------------------------------------------------------------------

type Switchable interface {
	// switchImpl replaces the internal metric implementation.
	// This allows converting from noop to real or real to noop without
	// breaking user references to the wrapper.
	switchImpl(newImpl any)
}

type SwitchableMetric interface {
	Switchable
	Metric
}

// baseSwitchableMetric provides common functionality for all switchable metrics.
//
// It holds a mutex and the current implementation of the metric.
// The following methods are provided to allow safe access to the internal metric.
// - switchImpl(newImpl any) to replace the internal implementation
// - IsNoop() bool to check if the current implementation is a noop
// - SetLevel(level Level) to set the level on the internal implementation
type baseSwitchableMetric[M Metric] struct {
	mu     sync.RWMutex
	impl   M
	isNoop bool
}

func newBaseSwitchableMetric[M Metric](impl M) *baseSwitchableMetric[M] {
	return &baseSwitchableMetric[M]{
		mu:   sync.RWMutex{},
		impl: impl,
	}
}

// switchImpl replaces the internal metric implementation.
//
// This is for internal use only, as it can break type safety if misused,
// (intentionally no type assertion check on newImpl)
func (b *baseSwitchableMetric[M]) switchImpl(newImpl any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.impl = newImpl.(M)
}

func (b *baseSwitchableMetric[M]) SetLevel(level Level) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.impl.SetLevel(level)
}

func (b *baseSwitchableMetric[M]) Name() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.impl.Name()
}

func (b *baseSwitchableMetric[M]) Help() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.impl.Help()
}

func (b *baseSwitchableMetric[M]) Type() MetricType {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.impl.Type()
}

func (b *baseSwitchableMetric[M]) Level() Level {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.impl.Level()
}

//--------------------------------------------------------------------------------
// Switchable Prime Metrics
//--------------------------------------------------------------------------------

// switchableCounter wraps a [Counter] implementation that can be switched
type switchableCounter struct {
	*baseSwitchableMetric[Counter]
}

func newSwitchableCounter(impl Counter, opts CounterOpts) *switchableCounter {
	return &switchableCounter{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCounter) Inc(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Inc(ctx)
}

func (s *switchableCounter) Add(ctx Context, value float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Add(ctx, value)
}

// switchableCounterVec wraps a [CounterVec] implementation that can be switched
type switchableCounterVec struct {
	*baseSwitchableMetric[CounterVec]
}

func newSwitchableCounterVec(impl CounterVec, opts CounterVecOpts) *switchableCounterVec {
	return &switchableCounterVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCounterVec) Inc(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Inc(ctx, labels)
}

func (s *switchableCounterVec) Add(ctx Context, value float64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Add(ctx, value, labels)
}

// switchableGauge wraps a [Gauge] implementation that can be switched
type switchableGauge struct {
	*baseSwitchableMetric[Gauge]
}

func newSwitchableGauge(impl Gauge, opts GaugeOpts) *switchableGauge {
	return &switchableGauge{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableGauge) Set(ctx Context, value float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Set(ctx, value)
}

func (s *switchableGauge) Inc(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Inc(ctx)
}

func (s *switchableGauge) Dec(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Dec(ctx)
}

func (s *switchableGauge) Add(ctx Context, value float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Add(ctx, value)
}

// switchableGaugeVec wraps a [GaugeVec] implementation that can be switched
type switchableGaugeVec struct {
	*baseSwitchableMetric[GaugeVec]
}

func newSwitchableGaugeVec(impl GaugeVec, opts GaugeVecOpts) *switchableGaugeVec {
	return &switchableGaugeVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableGaugeVec) Set(ctx Context, value float64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Set(ctx, value, labels)
}

func (s *switchableGaugeVec) Inc(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Inc(ctx, labels)
}

func (s *switchableGaugeVec) Dec(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Dec(ctx, labels)
}

func (s *switchableGaugeVec) Add(ctx Context, value float64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Add(ctx, value, labels)
}

type switchableHistogram struct {
	*baseSwitchableMetric[Histogram]
}

func newSwitchableHistogram(impl Histogram, opts HistogramOpts) *switchableHistogram {
	return &switchableHistogram{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableHistogram) Observe(ctx Context, value float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Observe(ctx, value)
}

type switchableHistogramVec struct {
	*baseSwitchableMetric[HistogramVec]
}

func newSwitchableHistogramVec(impl HistogramVec, opts HistogramVecOpts) *switchableHistogramVec {
	return &switchableHistogramVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableHistogramVec) Observe(ctx Context, value float64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Observe(ctx, value, labels)
}

type switchableSummary struct {
	*baseSwitchableMetric[Summary]
}

func newSwitchableSummary(impl Summary, opts SummaryOpts) *switchableSummary {
	return &switchableSummary{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableSummary) Observe(ctx Context, value float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Observe(ctx, value)
}

func (s *switchableSummary) Quantile(ctx Context, q float64) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Quantile(ctx, q)
}

type switchableSummaryVec struct {
	*baseSwitchableMetric[SummaryVec]
}

func newSwitchableSummaryVec(impl SummaryVec, opts SummaryVecOpts) *switchableSummaryVec {
	return &switchableSummaryVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableSummaryVec) Observe(ctx Context, value float64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Observe(ctx, value, labels)
}

func (s *switchableSummaryVec) Quantile(ctx Context, q float64, labels VecLabels) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Quantile(ctx, q, labels)
}

//--------------------------------------------------------------------------------
// Switchable Composite Metrics
//--------------------------------------------------------------------------------

// switchableTimer wraps a [Timer] implementation that can be switched
type switchableTimer struct {
	*baseSwitchableMetric[Timer]
}

func newSwitchableTimer(impl Timer, opts TimerOpts) *switchableTimer {
	return &switchableTimer{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableTimer) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableTimer) Start(ctx Context) func() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Start(ctx)
}

func (s *switchableTimer) Record(ctx Context, duration time.Duration) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Record(ctx, duration)
}

type switchableTimerVec struct {
	*baseSwitchableMetric[TimerVec]
}

func newSwitchableTimerVec(impl TimerVec, opts TimerVecOpts) *switchableTimerVec {
	return &switchableTimerVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableTimerVec) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableTimerVec) Start(ctx Context, labels VecLabels) func() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Start(ctx, labels)
}

func (s *switchableTimerVec) Record(ctx Context, duration time.Duration, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Record(ctx, duration, labels)
}

// switchableCache wraps a [Cache] implementation that can be switched
type switchableCache struct {
	*baseSwitchableMetric[Cache]
}

func newSwitchableCache(impl Cache, opts CacheOpts) *switchableCache {
	return &switchableCache{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCache) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableCache) Hit(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Hit(ctx)
}

func (s *switchableCache) Miss(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Miss(ctx)
}

func (s *switchableCache) SetSize(ctx Context, bytes int64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetSize(ctx, bytes)
}

type switchableCacheVec struct {
	*baseSwitchableMetric[CacheVec]
}

func newSwitchableCacheVec(impl CacheVec, opts CacheVecOpts) *switchableCacheVec {
	return &switchableCacheVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCacheVec) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableCacheVec) Hit(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Hit(ctx, labels)
}

func (s *switchableCacheVec) Miss(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Miss(ctx, labels)
}

func (s *switchableCacheVec) SetSize(ctx Context, bytes int64, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetSize(ctx, bytes, labels)
}

type switchablePool struct {
	*baseSwitchableMetric[Pool]
}

func newSwitchablePool(impl Pool, opts PoolOpts) *switchablePool {
	return &switchablePool{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchablePool) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchablePool) SetActive(ctx Context, count int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetActive(ctx, count)
}

func (s *switchablePool) SetIdle(ctx Context, count int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetIdle(ctx, count)
}

func (s *switchablePool) Acquired(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Acquired(ctx)
}

func (s *switchablePool) Released(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Released(ctx)
}

type switchablePoolVec struct {
	*baseSwitchableMetric[PoolVec]
}

func newSwitchablePoolVec(impl PoolVec, opts PoolVecOpts) *switchablePoolVec {
	return &switchablePoolVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchablePoolVec) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchablePoolVec) SetActive(ctx Context, count int, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetActive(ctx, count, labels)
}

func (s *switchablePoolVec) SetIdle(ctx Context, count int, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetIdle(ctx, count, labels)
}

func (s *switchablePoolVec) Acquired(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Acquired(ctx, labels)
}

func (s *switchablePoolVec) Released(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Released(ctx, labels)
}

type switchableCircuitBreaker struct {
	*baseSwitchableMetric[CircuitBreaker]
}

func newSwitchableCircuitBreaker(impl CircuitBreaker, opts CircuitBreakerOpts) *switchableCircuitBreaker {
	return &switchableCircuitBreaker{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCircuitBreaker) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableCircuitBreaker) SetState(ctx Context, state CircuitBreakerState) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetState(ctx, state)
}

func (s *switchableCircuitBreaker) Success(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Success(ctx)
}

func (s *switchableCircuitBreaker) Failure(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Failure(ctx)
}

type switchableCircuitBreakerVec struct {
	*baseSwitchableMetric[CircuitBreakerVec]
}

func newSwitchableCircuitBreakerVec(impl CircuitBreakerVec, opts CircuitBreakerVecOpts) *switchableCircuitBreakerVec {
	return &switchableCircuitBreakerVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableCircuitBreakerVec) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableCircuitBreakerVec) SetState(ctx Context, state CircuitBreakerState, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetState(ctx, state, labels)
}

func (s *switchableCircuitBreakerVec) Success(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Success(ctx, labels)
}

func (s *switchableCircuitBreakerVec) Failure(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Failure(ctx, labels)
}

type switchableQueue struct {
	*baseSwitchableMetric[Queue]
}

func newSwitchableQueue(impl Queue, opts QueueOpts) *switchableQueue {
	return &switchableQueue{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableQueue) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableQueue) SetDepth(ctx Context, depth int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetDepth(ctx, depth)
}

func (s *switchableQueue) Enqueued(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Enqueued(ctx)
}

func (s *switchableQueue) Dequeued(ctx Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Dequeued(ctx)
}

func (s *switchableQueue) SetWaitTime(ctx Context, duration time.Duration) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetWaitTime(ctx, duration)
}

type switchableQueueVec struct {
	*baseSwitchableMetric[QueueVec]
}

func newSwitchableQueueVec(impl QueueVec, opts QueueVecOpts) *switchableQueueVec {
	return &switchableQueueVec{
		baseSwitchableMetric: newBaseSwitchableMetric(impl),
	}
}

func (s *switchableQueueVec) Components() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Components()
}

func (s *switchableQueueVec) SetDepth(ctx Context, depth int, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetDepth(ctx, depth, labels)
}

func (s *switchableQueueVec) Enqueued(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Enqueued(ctx, labels)
}

func (s *switchableQueueVec) Dequeued(ctx Context, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.Dequeued(ctx, labels)
}

func (s *switchableQueueVec) SetWaitTime(ctx Context, duration time.Duration, labels VecLabels) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.impl.SetWaitTime(ctx, duration, labels)
}

var (
	__ctc_switchableCounter              Metric          = switchableCounter{}
	__ctc_switchableCounterPtr           Metric          = &switchableCounter{}
	__ctc_switchableCounterVec           Metric          = switchableCounterVec{}
	__ctc_switchableCounterVecPtr        Metric          = &switchableCounterVec{}
	__ctc_switchableGauge                Metric          = switchableGauge{}
	__ctc_switchableGaugePtr             Metric          = &switchableGauge{}
	__ctc_switchableGaugeVec             Metric          = switchableGaugeVec{}
	__ctc_switchableGaugeVecPtr          Metric          = &switchableGaugeVec{}
	__ctc_switchableHistogram            Metric          = switchableHistogram{}
	__ctc_switchableHistogramPtr         Metric          = &switchableHistogram{}
	__ctc_switchableHistogramVec         Metric          = switchableHistogramVec{}
	__ctc_switchableHistogramVecPtr      Metric          = &switchableHistogramVec{}
	__ctc_switchableSummary              Metric          = switchableSummary{}
	__ctc_switchableSummaryPtr           Metric          = &switchableSummary{}
	__ctc_switchableSummaryVec           Metric          = switchableSummaryVec{}
	__ctc_switchableSummaryVecPtr        Metric          = &switchableSummaryVec{}
	__ctc_switchableTimerPtr             CompositeMetric = &switchableTimer{}
	__ctc_switchableTimerVecPtr          CompositeMetric = &switchableTimerVec{}
	__ctc_switchableCachePtr             CompositeMetric = &switchableCache{}
	__ctc_switchableCacheVecPtr          CompositeMetric = &switchableCacheVec{}
	__ctc_switchablePoolPtr              CompositeMetric = &switchablePool{}
	__ctc_switchablePoolVecPtr           CompositeMetric = &switchablePoolVec{}
	__ctc_switchableCircuitBreakerPtr    CompositeMetric = &switchableCircuitBreaker{}
	__ctc_switchableCircuitBreakerVecPtr CompositeMetric = &switchableCircuitBreakerVec{}
	__ctc_switchableQueuePtr             CompositeMetric = &switchableQueue{}
	__ctc_switchableQueueVecPtr          CompositeMetric = &switchableQueueVec{}
)
