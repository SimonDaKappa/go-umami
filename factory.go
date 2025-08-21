package umami

import "time"

// factory implements the Factory interface
type factory struct {
	backend Backend
	group   string
	level   Level
	mask    MetricMask
}

// newFactory creates a new factory
func newFactory(backend Backend, group string, level Level, mask MetricMask) Factory {
	return &factory{
		backend: backend,
		group:   group,
		level:   level,
		mask:    mask,
	}
}

// Counter creates a counter with the given level and mask
func (f *factory) Counter(name string, level Level, mask MetricMask, labels ...string) Counter {
	fullName := f.group + "_" + name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopCounter{}
	}

	backend := f.backend.Counter(fullName, labels...)
	return &counter{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Gauge creates a gauge with the given level and mask
func (f *factory) Gauge(name string, level Level, mask MetricMask, labels ...string) Gauge {
	fullName := f.group + "_" + name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopGauge{}
	}

	backend := f.backend.Gauge(fullName, labels...)
	return &gauge{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Histogram creates a histogram with the given level and mask
func (f *factory) Histogram(name string, level Level, mask MetricMask, buckets []float64, labels ...string) Histogram {
	fullName := f.group + "_" + name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopHistogram{}
	}

	backend := f.backend.Histogram(fullName, buckets, labels...)
	return &histogram{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Summary creates a summary with the given level and mask
func (f *factory) Summary(name string, level Level, mask MetricMask, objectives map[float64]float64, labels ...string) Summary {
	fullName := f.group + "_" + name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopSummary{}
	}

	backend := f.backend.Summary(fullName, objectives, labels...)
	return &summary{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Timer creates a timer with the given level and mask
func (f *factory) Timer(name string, level Level, mask MetricMask, labels ...string) Timer {
	// Timer is built on top of histogram
	hist := f.Histogram(name+"_duration_seconds", level, mask, []float64{
		0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
	}, labels...)

	return &timer{histogram: hist}
}

// Cache creates cache metrics with the given level and mask
func (f *factory) Cache(name string, level Level, mask MetricMask, labels ...string) Cache {
	hits := f.Counter(name+"_hits", level, mask, labels...)
	misses := f.Counter(name+"_misses", level, mask, labels...)
	size := f.Gauge(name+"_size_bytes", level, mask, labels...)

	return &cache{
		hits:   hits,
		misses: misses,
		size:   size,
	}
}

// Pool creates pool metrics with the given level and mask
func (f *factory) Pool(name string, level Level, mask MetricMask, labels ...string) Pool {
	active := f.Gauge(name+"_active", level, mask, labels...)
	idle := f.Gauge(name+"_idle", level, mask, labels...)
	acquired := f.Counter(name+"_acquired", level, mask, labels...)
	released := f.Counter(name+"_released", level, mask, labels...)

	return &pool{
		active:   active,
		idle:     idle,
		acquired: acquired,
		released: released,
	}
}

// CircuitBreaker creates circuit breaker metrics with the given level and mask
func (f *factory) CircuitBreaker(name string, level Level, mask MetricMask, labels ...string) CircuitBreaker {
	state := f.Gauge(name+"_state", level, mask, labels...)
	successes := f.Counter(name+"_successes", level, mask, labels...)
	failures := f.Counter(name+"_failures", level, mask, labels...)

	return &circuitBreaker{
		state:     state,
		successes: successes,
		failures:  failures,
	}
}

// Queue creates queue metrics with the given level and mask
func (f *factory) Queue(name string, level Level, mask MetricMask, labels ...string) Queue {
	depth := f.Gauge(name+"_depth", level, mask, labels...)
	enqueued := f.Counter(name+"_enqueued", level, mask, labels...)
	dequeued := f.Counter(name+"_dequeued", level, mask, labels...)
	waitTime := f.Histogram(name+"_wait_time_seconds", level, mask, []float64{
		0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
	}, labels...)

	return &queue{
		depth:    depth,
		enqueued: enqueued,
		dequeued: dequeued,
		waitTime: waitTime,
	}
}

// Concrete metric implementations

// counter wraps a CounterBackend and implements early return
type counter struct {
	backend CounterBackend
	level   Level
	mask    MetricMask
}

func (c *counter) Inc(ctx Context) error {
	if !ctx.Enabled(c.level) || !ctx.EnabledMask(c.mask) {
		return nil // Early return
	}
	return c.backend.Inc()
}

func (c *counter) Add(ctx Context, value float64) error {
	if !ctx.Enabled(c.level) || !ctx.EnabledMask(c.mask) {
		return nil // Early return
	}
	return c.backend.Add(value)
}

// gauge wraps a GaugeBackend and implements early return
type gauge struct {
	backend GaugeBackend
	level   Level
	mask    MetricMask
}

func (g *gauge) Set(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Set(value)
}

func (g *gauge) Inc(ctx Context) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Inc()
}

func (g *gauge) Dec(ctx Context) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Dec()
}

func (g *gauge) Add(ctx Context, value float64) error {
	if !ctx.Enabled(g.level) || !ctx.EnabledMask(g.mask) {
		return nil // Early return
	}
	return g.backend.Add(value)
}

// histogram wraps a HistogramBackend and implements early return
type histogram struct {
	backend HistogramBackend
	level   Level
	mask    MetricMask
}

func (h *histogram) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(h.level) || !ctx.EnabledMask(h.mask) {
		return nil // Early return
	}
	return h.backend.Observe(value)
}

func (h *histogram) Time(ctx Context, fn func() error) error {
	if !ctx.Enabled(h.level) || !ctx.EnabledMask(h.mask) {
		return fn() // Execute function but don't time it
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)
	h.backend.Observe(duration.Seconds())
	return err
}

// summary wraps a SummaryBackend and implements early return
type summary struct {
	backend SummaryBackend
	level   Level
	mask    MetricMask
}

func (s *summary) Observe(ctx Context, value float64) error {
	if !ctx.Enabled(s.level) || !ctx.EnabledMask(s.mask) {
		return nil // Early return
	}
	return s.backend.Observe(value)
}

func (s *summary) Quantile(ctx Context, q float64) (float64, error) {
	if !ctx.Enabled(s.level) || !ctx.EnabledMask(s.mask) {
		return 0, nil // Early return with safe default
	}
	return s.backend.Quantile(q)
}

// Composite metric implementations

// timer combines a histogram for timing
type timer struct {
	histogram Histogram
}

func (t *timer) Start(ctx Context) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		t.histogram.Observe(ctx, duration.Seconds())
	}
}

func (t *timer) Record(ctx Context, duration time.Duration) error {
	return t.histogram.Observe(ctx, duration.Seconds())
}

// cache combines multiple metrics for cache operations
type cache struct {
	hits   Counter
	misses Counter
	size   Gauge
}

func (c *cache) Hit(ctx Context) error {
	return c.hits.Inc(ctx)
}

func (c *cache) Miss(ctx Context) error {
	return c.misses.Inc(ctx)
}

func (c *cache) SetSize(ctx Context, bytes int64) error {
	return c.size.Set(ctx, float64(bytes))
}

// pool combines multiple metrics for pool operations
type pool struct {
	active   Gauge
	idle     Gauge
	acquired Counter
	released Counter
}

func (p *pool) SetActive(ctx Context, count int) error {
	return p.active.Set(ctx, float64(count))
}

func (p *pool) SetIdle(ctx Context, count int) error {
	return p.idle.Set(ctx, float64(count))
}

func (p *pool) Acquired(ctx Context) error {
	return p.acquired.Inc(ctx)
}

func (p *pool) Released(ctx Context) error {
	return p.released.Inc(ctx)
}

// circuitBreaker combines multiple metrics for circuit breaker operations
type circuitBreaker struct {
	state     Gauge
	successes Counter
	failures  Counter
}

func (cb *circuitBreaker) SetState(ctx Context, state string) error {
	var value float64
	switch state {
	case "closed":
		value = 0
	case "open":
		value = 1
	case "half-open":
		value = 2
	default:
		value = -1
	}
	return cb.state.Set(ctx, value)
}

func (cb *circuitBreaker) Success(ctx Context) error {
	return cb.successes.Inc(ctx)
}

func (cb *circuitBreaker) Failure(ctx Context) error {
	return cb.failures.Inc(ctx)
}

// queue combines multiple metrics for queue operations
type queue struct {
	depth    Gauge
	enqueued Counter
	dequeued Counter
	waitTime Histogram
}

func (q *queue) SetDepth(ctx Context, depth int) error {
	return q.depth.Set(ctx, float64(depth))
}

func (q *queue) Enqueued(ctx Context) error {
	return q.enqueued.Inc(ctx)
}

func (q *queue) Dequeued(ctx Context) error {
	return q.dequeued.Inc(ctx)
}

func (q *queue) SetWaitTime(ctx Context, duration time.Duration) error {
	return q.waitTime.Observe(ctx, duration.Seconds())
}
