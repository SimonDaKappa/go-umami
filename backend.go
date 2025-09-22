package umami

//--------------------------------------------------------------------------------
// File: backend.go
//
// This file contains the definition of the [Backend] interface and related
// adapter interfaces for the umami metrics library.
//--------------------------------------------------------------------------------

// Backend defines the interface that concrete backend implementations must
// satisfy.
//
// A backend is responsible for creating and managing metric
// instances. Each metric instance is represented by an adapter interface
// (e.g. CounterAdapter, GaugeAdapter, etc.) that provides methods for
// manipulating the metric's value.
//
// Backends must minimally support the below defined adapter interfaces,
// but may choose to implement additional functionality as needed. That is,
// a backend must provide an adapter for each Basic [Metric] type, but may
// choose to implement additional metric types or features.
type Backend interface {
	Counter(opts CounterOpts) CounterAdapter
	CounterVec(opts CounterVecOpts) CounterVecAdapter
	Gauge(opts GaugeOpts) GaugeAdapter
	GaugeVec(opts GaugeVecOpts) GaugeVecAdapter
	Histogram(opts HistogramOpts) HistogramAdapter
	HistogramVec(opts HistogramVecOpts) HistogramVecAdapter
	Summary(opts SummaryOpts) SummaryAdapter
	SummaryVec(opts SummaryVecOpts) SummaryVecAdapater
	Name() string
}

// CounterAdapter defines the interface for counter metrics that concrete
// backend adapter implementations must satisfy.
type CounterAdapter interface {
	Inc() error
	Add(value float64) error
}

// CounterVecAdapter defines the interface for partitioned counter metrics that
// concrete backend adapter implementations must satisfy.
type CounterVecAdapter interface {
	Inc(labels VecLabels) error
	Add(value float64, labels VecLabels) error
}

// GaugeAdapter defines the interface for gauge metrics that concrete
// backend adapter implementations must satisfy.
type GaugeAdapter interface {
	Set(value float64) error
	Inc() error
	Dec() error
	Add(value float64) error
}

type GaugeVecAdapter interface {
	Set(value float64, labels VecLabels) error
	Inc(labels VecLabels) error
	Dec(labels VecLabels) error
	Add(value float64, labels VecLabels) error
}

// HistogramAdapter defines the interface for histogram metrics that concrete
// backend adapter implementations must satisfy.
type HistogramAdapter interface {
	Observe(value float64) error
}

type HistogramVecAdapter interface {
	Observe(value float64, labels VecLabels) error
}

// SummaryAdapter defines the interface for summary metrics that concrete
// backend adapter implementations must satisfy.
type SummaryAdapter interface {
	Observe(value float64) error
	Quantile(q float64) (float64, error)
}

type SummaryVecAdapater interface {
	Observe(value float64, labels VecLabels) error
	Quantile(q float64, labels VecLabels) (float64, error)
}

const (
	BackendNoneName string = "none"
)
