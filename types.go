package umami

import "time"

// VecLabels is a type that represents a set partition keys to values
type VecLabels map[string]string

// CounterOpts contains options for counter metrics
type CounterOpts struct {
	Name string
	Help string
}

// baseCounter wraps a CounterBackend and implements early return
type baseCounter struct {
	backend CounterBackend
	level   Level
	mask    Mask
}

func (c *baseCounter) Inc(ctx Context) error {
	if !ctx.Enabled(c.level) || !ctx.EnabledMask(c.mask) {
		return nil // Early return
	}
	return c.backend.Inc()
}

func (c *baseCounter) Add(ctx Context, value float64) error {
	if !ctx.Enabled(c.level) || !ctx.EnabledMask(c.mask) {
		return nil // Early return
	}
	return c.backend.Add(value)
}

type CounterVecOpts struct {
	Name   string
	Help   string
	Labels []string
}

type baseCounterVec struct {
	backend CounterVecBackend
	level   Level
	mask    Mask
}

func (cv *baseCounterVec) Inc(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(cv.level) || !ctx.EnabledMask(cv.mask) {
		return nil // Early return
	}
	return cv.backend.Inc(labels)
}

func (cv *baseCounterVec) Add(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(cv.level) || !ctx.EnabledMask(cv.mask) {
		return nil // Early return
	}
	return cv.backend.Add(value, labels)
}

// GaugeOpts can be extended by backends for additional options
type GaugeOpts struct {
	Name string
	Help string
}

// baseGauge wraps a GaugeBackend and implements early return
type baseGauge struct {
	backend GaugeBackend
	level   Level
	mask    Mask
}

func (g *baseGauge) Set(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Set(value)
}

func (g *baseGauge) Inc(ctx Context) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Inc()
}

func (g *baseGauge) Dec(ctx Context) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Dec()
}

func (g *baseGauge) Add(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Add(value)
}

type GaugeVecOpts struct {
	Name   string
	Help   string
	Labels []string
}

// gauge wraps a GaugeBackend and implements early return
type baseGaugeVec struct {
	backend GaugeVecBackend
	level   Level
	mask    Mask
}

func (gv *baseGaugeVec) Set(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(gv.level) || !ctx.EnabledMask(gv.mask) {
		return nil // Early return
	}
	return gv.backend.Set(value, labels)
}

func (gv *baseGaugeVec) Inc(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(gv.level) || !ctx.EnabledMask(gv.mask) {
		return nil // Early return
	}
	return gv.backend.Inc(labels)
}

func (gv *baseGaugeVec) Dec(ctx Context, labels VecLabels) error {
	if !ctx.Enabled(gv.level) || !ctx.EnabledMask(gv.mask) {
		return nil // Early return
	}
	return gv.backend.Dec(labels)
}

func (gv *baseGaugeVec) Add(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(gv.level) || !ctx.EnabledMask(gv.mask) {
		return nil // Early return
	}
	return gv.backend.Add(value, labels)
}

// HistogramOpts can be extended by backends for additional options
type HistogramOpts struct {
	Name    string
	Help    string
	Buckets []float64
}

// baseHistogram wraps a HistogramBackend and implements early return
type baseHistogram struct {
	backend HistogramBackend
	level   Level
	mask    Mask
}

func (h *baseHistogram) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(h.level) || !ctx.EnabledMask(h.mask) {
		return nil // Early return
	}
	return h.backend.Observe(value)
}

func (h *baseHistogram) Time(ctx Context, fn func() error) error {
	if !ctx.Enabled(h.level) || !ctx.EnabledMask(h.mask) {
		return fn() // Execute function but don't time it
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)
	h.backend.Observe(duration.Seconds())
	return err
}

type HistogramVecOpts struct {
	Name    string
	Help    string
	Labels  []string
	Buckets []float64
}

// histogram wraps a HistogramBackend and implements early return
type baseHistogramVec struct {
	backend HistogramVecBackend
	level   Level
	mask    Mask
}

func (hv *baseHistogramVec) Observe(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(hv.level) || !ctx.EnabledMask(hv.mask) {
		return nil // Early return
	}
	return hv.backend.Observe(value, labels)
}

func (hv *baseHistogramVec) Time(ctx Context, fn func() error, labels VecLabels) error {
	if !ctx.Enabled(hv.level) || !ctx.EnabledMask(hv.mask) {
		return fn() // Execute function but don't time it
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)
	hv.backend.Observe(duration.Seconds(), labels)
	return err
}

type TimerOpts struct {
	HistOpts HistogramOpts
}

// baseTimer combines a histogram for timing
type baseTimer struct {
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

type TimerVecOpts struct {
	HistVecOpts HistogramVecOpts
}

// timer combines a histogram for timing
type baseTimerVec struct {
	histogram HistogramVec
}

func (tv *baseTimerVec) Start(ctx Context, labels VecLabels) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		tv.histogram.Observe(ctx, duration.Seconds(), labels)
	}
}

func (tv *baseTimerVec) Record(ctx Context, duration time.Duration, labels VecLabels) error {
	return tv.histogram.Observe(ctx, duration.Seconds(), labels)
}

type CacheOpts struct {
	HitOpts  CounterOpts
	MissOpts CounterOpts
	SizeOpts GaugeOpts
}

// baseCache combines multiple metrics for baseCache operations
type baseCache struct {
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

type CacheVecOpts struct {
	HitVecOpts  CounterVecOpts
	MissVecOpts CounterVecOpts
	SizeVecOpts GaugeVecOpts
}

// cache combines multiple metrics for cache operations
type baseCacheVec struct {
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

type PoolOpts struct {
	ActiveOpts   GaugeOpts
	IdleOpts     GaugeOpts
	AcquiredOpts CounterOpts
	ReleasedOpts CounterOpts
}

// basePool combines multiple metrics for basePool operations
type basePool struct {
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

type PoolVecOpts struct {
	ActiveVecOpts   GaugeVecOpts
	IdleVecOpts     GaugeVecOpts
	AcquiredVecOpts CounterVecOpts
	ReleasedVecOpts CounterVecOpts
}

// pool combines multiple metrics for pool operations
type basePoolVec struct {
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

type CircuitBreakerOpts struct {
	StateOpts   GaugeOpts
	SuccessOpts CounterOpts
	FailureOpts CounterOpts
}

// baseCircuitBreaker combines multiple metrics for circuit breaker operations
type baseCircuitBreaker struct {
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

type CircuitBreakerVecOpts struct {
	StateVecOpts   GaugeVecOpts
	SuccessVecOpts CounterVecOpts
	FailureVecOpts CounterVecOpts
}

type CircuitBreakerState uint8

const (
	CircuitBreakerStateClosed   CircuitBreakerState = 0
	CircuitBreakerStateOpen     CircuitBreakerState = 1
	CircuitBreakerStateHalfOpen CircuitBreakerState = 2
)

// circuitBreaker combines multiple metrics for circuit breaker operations
type baseCircuitBreakerVec struct {
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

type SummaryOpts struct {
	Name       string
	Help       string
	Objectives map[float64]float64
}

// baseSummary wraps a SummaryBackend and implements early return
type baseSummary struct {
	backend SummaryBackend
	level   Level
	mask    Mask
}

func (s *baseSummary) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(s.level) || !ctx.EnabledMask(s.mask) {
		return nil // Early return
	}
	return s.backend.Observe(value)
}

func (s *baseSummary) Quantile(ctx Context, q float64) (float64, error) {
	if !ctx.Enabled(s.level) || !ctx.EnabledMask(s.mask) {
		return 0, nil // Early return with safe default
	}
	return s.backend.Quantile(q)
}

type SummaryVecOpts struct {
	Name       string
	Help       string
	Labels     []string
	Objectives map[float64]float64
}

// summary wraps a SummaryBackend and implements early return
type baseSummaryVec struct {
	backend SummaryVecBackend
	level   Level
	mask    Mask
}

func (sv *baseSummaryVec) Observe(ctx Context, value float64, labels VecLabels) error {
	if !ctx.Enabled(sv.level) || !ctx.EnabledMask(sv.mask) {
		return nil // Early return
	}
	return sv.backend.Observe(value, labels)
}

func (sv *baseSummaryVec) Quantile(ctx Context, q float64, labels VecLabels) (float64, error) {
	if !ctx.Enabled(sv.level) || !ctx.EnabledMask(sv.mask) {
		return 0, nil // Early return with safe default
	}
	return sv.backend.Quantile(q, labels)
}

type QueueOpts struct {
	DepthOpts    GaugeOpts
	EnqueuedOpts CounterOpts
	DequeuedOpts CounterOpts
	WaitTimeOpts HistogramOpts
}

// baseQueue combines multiple metrics for baseQueue operations
type baseQueue struct {
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

type QueueVecOpts struct {
	DepthVecOpts    GaugeVecOpts
	EnqueuedVecOpts CounterVecOpts
	DequeuedVecOpts CounterVecOpts
	WaitTimeVecOpts HistogramVecOpts
}

// queue combines multiple metrics for queue operations
type baseQueueVec struct {
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
