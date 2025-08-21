package umami

import "time"

// No-op implementations for disabled metrics
// These provide zero-cost operations when metrics are disabled

// noopCounter implements Counter interface with no-op operations
type noopCounter struct{}

func (n *noopCounter) Inc(ctx Context) error                { return nil }
func (n *noopCounter) Add(ctx Context, value float64) error { return nil }

// noopGauge implements Gauge interface with no-op operations
type noopGauge struct{}

func (n *noopGauge) Set(ctx Context, value float64) error { return nil }
func (n *noopGauge) Inc(ctx Context) error                { return nil }
func (n *noopGauge) Dec(ctx Context) error                { return nil }
func (n *noopGauge) Add(ctx Context, value float64) error { return nil }

// noopHistogram implements Histogram interface with no-op operations
type noopHistogram struct{}

func (n *noopHistogram) Observe(ctx Context, value float64) error { return nil }
func (n *noopHistogram) Time(ctx Context, fn func() error) error  { return fn() }

// noopSummary implements Summary interface with no-op operations
type noopSummary struct{}

func (n *noopSummary) Observe(ctx Context, value float64) error         { return nil }
func (n *noopSummary) Quantile(ctx Context, q float64) (float64, error) { return 0, nil }

// noopTimer implements Timer interface with no-op operations
type noopTimer struct{}

func (n *noopTimer) Start(ctx Context) func()                         { return func() {} }
func (n *noopTimer) Record(ctx Context, duration time.Duration) error { return nil }

// noopCache implements Cache interface with no-op operations
type noopCache struct{}

func (n *noopCache) Hit(ctx Context) error                  { return nil }
func (n *noopCache) Miss(ctx Context) error                 { return nil }
func (n *noopCache) SetSize(ctx Context, bytes int64) error { return nil }

// noopPool implements Pool interface with no-op operations
type noopPool struct{}

func (n *noopPool) SetActive(ctx Context, count int) error { return nil }
func (n *noopPool) SetIdle(ctx Context, count int) error   { return nil }
func (n *noopPool) Acquired(ctx Context) error             { return nil }
func (n *noopPool) Released(ctx Context) error             { return nil }

// noopCircuitBreaker implements CircuitBreaker interface with no-op operations
type noopCircuitBreaker struct{}

func (n *noopCircuitBreaker) SetState(ctx Context, state string) error { return nil }
func (n *noopCircuitBreaker) Success(ctx Context) error                { return nil }
func (n *noopCircuitBreaker) Failure(ctx Context) error                { return nil }

// noopQueue implements Queue interface with no-op operations
type noopQueue struct{}

func (n *noopQueue) SetDepth(ctx Context, depth int) error                 { return nil }
func (n *noopQueue) Enqueued(ctx Context) error                            { return nil }
func (n *noopQueue) Dequeued(ctx Context) error                            { return nil }
func (n *noopQueue) SetWaitTime(ctx Context, duration time.Duration) error { return nil }
