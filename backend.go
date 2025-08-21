package umami

type Backend interface {
	Counter(opts CounterOpts) CounterBackend
	CounterVec(opts CounterVecOpts) CounterVecBackend
	Gauge(opts GaugeOpts) GaugeBackend
	GaugeVec(opts GaugeVecOpts) GaugeVecBackend
	Histogram(opts HistogramOpts) HistogramBackend
	HistogramVec(opts HistogramVecOpts) HistogramVecBackend
	Summary(opts SummaryOpts) SummaryBackend
	SummaryVec(opts SummaryVecOpts) SummaryVecBackend
	Name() string
}

const (
	BackendNoneName string = "none"
)

// CounterBackend defines the interface for counter metrics that concrete
// backend adapter implementations must satisfy.
type CounterBackend interface {
	Inc() error
	Add(value float64) error
}

// CounterVecBackend defines the interface for partitioned counter metrics that
// concrete backend adapter implementations must satisfy.
type CounterVecBackend interface {
	Inc(labels VecLabels) error
	Add(value float64, labels VecLabels) error
}

// GaugeBackend defines the interface for gauge metrics that concrete
// backend adapter implementations must satisfy.
type GaugeBackend interface {
	Set(value float64) error
	Inc() error
	Dec() error
	Add(value float64) error
}

type GaugeVecBackend interface {
	Set(value float64, labels VecLabels) error
	Inc(labels VecLabels) error
	Dec(labels VecLabels) error
	Add(value float64, labels VecLabels) error
}

// HistogramBackend defines the interface for histogram metrics that concrete
// backend adapter implementations must satisfy.
type HistogramBackend interface {
	Observe(value float64) error
}

type HistogramVecBackend interface {
	Observe(value float64, labels VecLabels) error
}

// SummaryBackend defines the interface for summary metrics that concrete
// backend adapter implementations must satisfy.
type SummaryBackend interface {
	Observe(value float64) error
	Quantile(q float64) (float64, error)
}

type SummaryVecBackend interface {
	Observe(value float64, labels VecLabels) error
	Quantile(q float64, labels VecLabels) (float64, error)
}
