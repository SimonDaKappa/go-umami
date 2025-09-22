package umami

//--------------------------------------------------------------------------------
// File: metric_types.go
//
// This file contains the metric type classification system used for noop
// conversion and metric management.
//--------------------------------------------------------------------------------

// MetricType represents the category/type of a metric
type MetricType uint8

const (
	MetricTypeBasic MetricType = iota
	MetricTypeComposite
	// Future metric types can be added here
)

// String returns a string representation of the MetricType
func (mt MetricType) String() string {
	switch mt {
	case MetricTypeBasic:
		return "Basic"
	case MetricTypeComposite:
		return "Composite"
	default:
		return "Unknown"
	}
}

// NoopMetricInfo stores the information needed to recreate a noop metric
// as a real implementation
type NoopMetricInfo struct {
	Type    MetricType
	Options interface{} // Will store the original opts (CounterOpts, GaugeOpts, etc.)
	Level   Level       // The level the metric was created with
}
