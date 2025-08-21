package umami

import "strings"

// Level represents the importance/verbosity level of a metric
// Similar to log levels in slog/zap
type Level int8

const (
	// LevelDisabled disables all metrics
	LevelDisabled Level = iota - 1

	// LevelCritical - Essential business metrics (SLA, errors, core counters)
	// Always enabled in production
	LevelCritical

	// LevelImportant - Important operational metrics (latency, throughput)
	// Usually enabled in production
	LevelImportant

	// LevelDebug - Debugging metrics (detailed breakdowns, internal state)
	// Often disabled in production
	LevelDebug

	// LevelVerbose - Very detailed metrics (per-user, per-request)
	// Usually only enabled for troubleshooting
	LevelVerbose
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelDisabled:
		return "DISABLED"
	case LevelCritical:
		return "CRITICAL"
	case LevelImportant:
		return "IMPORTANT"
	case LevelDebug:
		return "DEBUG"
	case LevelVerbose:
		return "VERBOSE"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a level string into a Level
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DISABLED":
		return LevelDisabled
	case "CRITICAL":
		return LevelCritical
	case "IMPORTANT":
		return LevelImportant
	case "DEBUG":
		return LevelDebug
	case "VERBOSE":
		return LevelVerbose
	default:
		return LevelImportant // Safe default
	}
}

// Enabled returns true if this level should be processed given the configured level
func (l Level) Enabled(configuredLevel Level) bool {
	return l <= configuredLevel && configuredLevel != LevelDisabled
}

// MetricMask allows fine-grained control over which metrics are enabled
// Similar to log masks but for metrics
type MetricMask uint64

const (
	// Core business metrics
	MaskCounters   MetricMask = 1 << iota // Basic counters
	MaskLatency                           // Latency/timing metrics
	MaskThroughput                        // Rate/throughput metrics
	MaskErrors                            // Error metrics

	// Operational metrics
	MaskResources   // CPU, memory, disk
	MaskQueues      // Queue depths, processing
	MaskConnections // DB connections, pools
	MaskCache       // Cache hit/miss rates

	// Advanced metrics
	MaskCircuitBreaker // Circuit breaker states
	MaskHealth         // Health check results
	MaskSecurity       // Auth, rate limiting
	MaskPerformance    // Detailed performance breakdowns

	// Development/Debug metrics
	MaskInternal   // Internal state metrics
	MaskPerUser    // Per-user metrics
	MaskPerRequest // Per-request metrics
	MaskDetailed   // Very detailed breakdowns

	// Convenience masks
	MaskNone       MetricMask = 0
	MaskEssential             = MaskCounters | MaskLatency | MaskErrors
	MaskProduction            = MaskEssential | MaskThroughput | MaskResources | MaskQueues
	MaskAll        MetricMask = ^MetricMask(0)
)

// Has returns true if the mask has the specified flag
func (m MetricMask) Has(flag MetricMask) bool {
	return m&flag != 0
}

// Add adds a flag to the mask
func (m MetricMask) Add(flag MetricMask) MetricMask {
	return m | flag
}

// Remove removes a flag from the mask
func (m MetricMask) Remove(flag MetricMask) MetricMask {
	return m &^ flag
}

// String returns a human-readable representation of the mask
func (m MetricMask) String() string {
	if m == MaskNone {
		return "NONE"
	}
	if m == MaskAll {
		return "ALL"
	}

	var parts []string
	flags := map[MetricMask]string{
		MaskCounters:       "COUNTERS",
		MaskLatency:        "LATENCY",
		MaskThroughput:     "THROUGHPUT",
		MaskErrors:         "ERRORS",
		MaskResources:      "RESOURCES",
		MaskQueues:         "QUEUES",
		MaskConnections:    "CONNECTIONS",
		MaskCache:          "CACHE",
		MaskCircuitBreaker: "CIRCUIT_BREAKER",
		MaskHealth:         "HEALTH",
		MaskSecurity:       "SECURITY",
		MaskPerformance:    "PERFORMANCE",
		MaskInternal:       "INTERNAL",
		MaskPerUser:        "PER_USER",
		MaskPerRequest:     "PER_REQUEST",
		MaskDetailed:       "DETAILED",
	}

	for flag, name := range flags {
		if m.Has(flag) {
			parts = append(parts, name)
		}
	}

	return strings.Join(parts, "|")
}
