package umami

// factory implements the [Factory] interface
type factory struct {
	backend Backend
	group   string
	level   Level
	mask    Mask
}

// newFactory creates a new factory
func newFactory(backend Backend, group string, level Level, mask Mask) Factory {
	return &factory{
		backend: backend,
		group:   group,
		level:   level,
		mask:    mask,
	}
}

// Counter creates a counter with the given level and mask
func (f *factory) Counter(opts CounterOpts, level Level, mask Mask) Counter {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopCounter{}
	}

	cbackend := f.backend.Counter(opts)
	return &baseCounter{
		backend: cbackend,
		level:   level,
		mask:    mask,
	}
}

func (f *factory) CounterVec(opts CounterVecOpts, level Level, mask Mask) CounterVec {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopCounterVec{}
	}

	cbackend := f.backend.CounterVec(opts)
	return &baseCounterVec{
		backend: cbackend,
		level:   level,
		mask:    mask,
	}
}

// Gauge creates a gauge with the given level and mask
func (f *factory) Gauge(opts GaugeOpts, level Level, mask Mask) Gauge {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopGauge{}
	}

	backend := f.backend.Gauge(opts)
	return &baseGauge{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// GaugeVec creates a gauge vector with the given level and mask
func (f *factory) GaugeVec(opts GaugeVecOpts, level Level, mask Mask) GaugeVec {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopGaugeVec{}
	}

	backend := f.backend.GaugeVec(opts)
	return &baseGaugeVec{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Histogram creates a histogram with the given level and mask
func (f *factory) Histogram(opts HistogramOpts, level Level, mask Mask) Histogram {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopHistogram{}
	}

	backend := f.backend.Histogram(opts)
	return &baseHistogram{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// HistogramVec creates a histogram vector with the given level and mask
func (f *factory) HistogramVec(opts HistogramVecOpts, level Level, mask Mask) HistogramVec {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopHistogramVec{}
	}

	backend := f.backend.HistogramVec(opts)
	return &baseHistogramVec{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Summary creates a summary with the given level and mask
func (f *factory) Summary(opts SummaryOpts, level Level, mask Mask) Summary {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopSummary{}
	}

	backend := f.backend.Summary(opts)
	return &baseSummary{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// SummaryVec creates a summary vector with the given level and mask
func (f *factory) SummaryVec(opts SummaryVecOpts, level Level, mask Mask) SummaryVec {
	opts.Name = f.group + "_" + opts.Name

	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopSummaryVec{}
	}

	backend := f.backend.SummaryVec(opts)
	return &baseSummaryVec{
		backend: backend,
		level:   level,
		mask:    mask,
	}
}

// Timer creates a timer with the given level and mask
func (f *factory) Timer(opts TimerOpts, level Level, mask Mask) Timer {
	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopTimer{}
	}

	// Timer is built on top of histogram
	hist := f.Histogram(opts.HistOpts, level, mask)

	return &baseTimer{histogram: hist}
}

// TimerVec creates a timer vector with the given level and mask
func (f *factory) TimerVec(opts TimerVecOpts, level Level, mask Mask) TimerVec {
	// Check if this metric should be enabled
	if !level.Enabled(f.level) || !f.mask.Has(mask) {
		return &noopTimerVec{}
	}

	// TimerVec is built on top of HistogramVec
	hist := f.HistogramVec(opts.HistVecOpts, level, mask)

	return &baseTimerVec{histogram: hist}
}

// Cache creates cache metrics with the given level and mask
func (f *factory) Cache(opts CacheOpts, level Level, mask Mask) Cache {
	hits := f.Counter(opts.HitOpts, level, mask)
	misses := f.Counter(opts.MissOpts, level, mask)
	size := f.Gauge(opts.SizeOpts, level, mask)

	return &baseCache{
		hits:   hits,
		misses: misses,
		size:   size,
	}
}

// CacheVec creates a cache vector with the given level and mask
func (f *factory) CacheVec(opts CacheVecOpts, level Level, mask Mask) CacheVec {
	hits := f.CounterVec(opts.HitVecOpts, level, mask)
	misses := f.CounterVec(opts.MissVecOpts, level, mask)
	size := f.GaugeVec(opts.SizeVecOpts, level, mask)

	return &baseCacheVec{
		hits:   hits,
		misses: misses,
		size:   size,
	}
}

// Pool creates pool metrics with the given level and mask
func (f *factory) Pool(opts PoolOpts, level Level, mask Mask) Pool {
	active := f.Gauge(opts.ActiveOpts, level, mask)
	idle := f.Gauge(opts.IdleOpts, level, mask)
	acquired := f.Counter(opts.AcquiredOpts, level, mask)
	released := f.Counter(opts.ReleasedOpts, level, mask)

	return &basePool{
		active:   active,
		idle:     idle,
		acquired: acquired,
		released: released,
	}
}

// PoolVec creates a pool vector with the given level and mask
func (f *factory) PoolVec(opts PoolVecOpts, level Level, mask Mask) PoolVec {
	active := f.GaugeVec(opts.ActiveVecOpts, level, mask)
	idle := f.GaugeVec(opts.IdleVecOpts, level, mask)
	acquired := f.CounterVec(opts.AcquiredVecOpts, level, mask)
	released := f.CounterVec(opts.ReleasedVecOpts, level, mask)

	return &basePoolVec{
		active:   active,
		idle:     idle,
		acquired: acquired,
		released: released,
	}
}

// CircuitBreaker creates circuit breaker metrics with the given level and mask
func (f *factory) CircuitBreaker(opts CircuitBreakerOpts, level Level, mask Mask) CircuitBreaker {
	state := f.Gauge(opts.StateOpts, level, mask)
	successes := f.Counter(opts.SuccessOpts, level, mask)
	failures := f.Counter(opts.FailureOpts, level, mask)

	return &baseCircuitBreaker{
		state:     state,
		successes: successes,
		failures:  failures,
	}
}

// CircuitBreakerVec creates a circuit breaker vector with the given level and mask
func (f *factory) CircuitBreakerVec(opts CircuitBreakerVecOpts, level Level, mask Mask) CircuitBreakerVec {
	state := f.GaugeVec(opts.StateVecOpts, level, mask)
	successes := f.CounterVec(opts.SuccessVecOpts, level, mask)
	failures := f.CounterVec(opts.FailureVecOpts, level, mask)

	return &baseCircuitBreakerVec{
		state:     state,
		successes: successes,
		failures:  failures,
	}
}

// Queue creates queue metrics with the given level and mask
func (f *factory) Queue(opts QueueOpts, level Level, mask Mask) Queue {
	depth := f.Gauge(opts.DepthOpts, level, mask)
	enqueued := f.Counter(opts.EnqueuedOpts, level, mask)
	dequeued := f.Counter(opts.DequeuedOpts, level, mask)
	waitTime := f.Histogram(opts.WaitTimeOpts, level, mask)

	return &baseQueue{
		depth:    depth,
		enqueued: enqueued,
		dequeued: dequeued,
		waitTime: waitTime,
	}
}

// QueueVec creates a queue vector with the given level and mask
func (f *factory) QueueVec(opts QueueVecOpts, level Level, mask Mask) QueueVec {
	depth := f.GaugeVec(opts.DepthVecOpts, level, mask)
	enqueued := f.CounterVec(opts.EnqueuedVecOpts, level, mask)
	dequeued := f.CounterVec(opts.DequeuedVecOpts, level, mask)
	waitTime := f.HistogramVec(opts.WaitTimeVecOpts, level, mask)

	return &baseQueueVec{
		depth:    depth,
		enqueued: enqueued,
		dequeued: dequeued,
		waitTime: waitTime,
	}
}
