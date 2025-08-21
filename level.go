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

var (
	LevelDisabledStr  = "DISABLED"
	LevelCriticalStr  = "CRITICAL"
	LevelImportantStr = "IMPORTANT"
	LevelDebugStr     = "DEBUG"
	LevelVerboseStr   = "VERBOSE"
	LevelUnknownStr   = "UNKNOWN"
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case LevelDisabled:
		return LevelDisabledStr
	case LevelCritical:
		return LevelCriticalStr
	case LevelImportant:
		return LevelImportantStr
	case LevelDebug:
		return LevelDebugStr
	case LevelVerbose:
		return LevelVerboseStr
	default:
		// TODO $$$SIMON Should be logged?
		return LevelUnknownStr
	}
}

// ParseLevel parses a level string into a Level
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case LevelDisabledStr:
		return LevelDisabled
	case LevelCriticalStr:
		return LevelCritical
	case LevelImportantStr:
		return LevelImportant
	case LevelDebugStr:
		return LevelDebug
	case LevelVerboseStr:
		return LevelVerbose
	default:
		// TODO $$$SIMON Should be logged?
		return LevelImportant // Safe default
	}
}

// Enabled returns true if this level should be processed given the configured level
func (l Level) Enabled(configuredLevel Level) bool {
	return l <= configuredLevel && configuredLevel != LevelDisabled
}

// Mask allows fine-grained control over which metrics are enabled
// Similar to log masks but for metrics
type Mask uint64

const (
	// Core business metrics
	MaskCounters   Mask = 1 << iota // Basic counters
	MaskLatency    Mask = 1 << 0x2  // Latency/timing metrics
	MaskThroughput Mask = 1 << 0x3  // Rate/throughput metrics
	MaskErrors     Mask = 1 << 0x4  // Error metrics
	// Operational metrics
	MaskResources   Mask = 1 << 0x5 // CPU, memory, disk
	MaskQueues      Mask = 1 << 0x6 // Queue depths, processing
	MaskConnections Mask = 1 << 0x7 // DB connections, pools
	MaskCache       Mask = 1 << 0x8 // Cache hit/miss rates
	// Advanced metrics
	MaskCircuitBreaker Mask = 1 << 0x9 // Circuit breaker states
	MaskHealth         Mask = 1 << 0xA // Health check results
	MaskSecurity       Mask = 1 << 0xB // Auth, rate limiting
	MaskPerformance    Mask = 1 << 0xC // Detailed performance breakdowns
	// Development/Debug metrics
	MaskInternal   Mask = 1 << 0xD  // Internal state metrics
	MaskPerUser    Mask = 1 << 0xE  // Per-user metrics
	MaskPerRequest Mask = 1 << 0xF  // Per-request metrics
	MaskDetailed   Mask = 1 << 0x10 // Very detailed breakdowns

	// Convenience masks
	MaskNone       Mask = 0
	MaskEssential  Mask = MaskCounters | MaskLatency | MaskErrors
	MaskProduction Mask = MaskEssential | MaskThroughput | MaskResources | MaskQueues
	MaskAll        Mask = ^Mask(0)
)

const (
	MaskCountersStr       string = "COUNTERS"
	MaskLatencyStr        string = "LATENCY"
	MaskThroughputStr     string = "THROUGHPUT"
	MaskErrorsStr         string = "ERRORS"
	MaskResourcesStr      string = "RESOURCES"
	MaskQueuesStr         string = "QUEUES"
	MaskConnectionsStr    string = "CONNECTIONS"
	MaskCacheStr          string = "CACHE"
	MaskCircuitBreakerStr string = "CIRCUIT_BREAKER"
	MaskHealthStr         string = "HEALTH"
	MaskSecurityStr       string = "SECURITY"
	MaskPerformanceStr    string = "PERFORMANCE"
	MaskInternalStr       string = "INTERNAL"
	MaskPerUserStr        string = "PER_USER"
	MaskPerRequestStr     string = "PER_REQUEST"
	MaskDetailedStr       string = "DETAILED"

	MaskNoneStr       string = "NONE"
	MaskEssentialStr  string = "ESSENTIAL"
	MaskProductionStr string = "PRODUCTION"
	MaskAllStr        string = "ALL"
)

// Has returns true if the mask has the specified flag
func (m Mask) Has(flag Mask) bool {
	return m&flag != 0
}

// Add adds a flag to the mask
func (m Mask) Add(flag Mask) Mask {
	return m | flag
}

// Remove removes a flag from the mask
func (m Mask) Remove(flag Mask) Mask {
	return m &^ flag
}

// String returns a human-readable representation of the mask
func (m Mask) String() string {
	if m == MaskNone {
		return MaskNoneStr
	}
	if m == MaskAll {
		return MaskAllStr
	}

	var parts []string
	flags := map[Mask]string{
		MaskCounters:       MaskCountersStr,
		MaskLatency:        MaskLatencyStr,
		MaskThroughput:     MaskThroughputStr,
		MaskErrors:         MaskErrorsStr,
		MaskResources:      MaskResourcesStr,
		MaskQueues:         MaskQueuesStr,
		MaskConnections:    MaskConnectionsStr,
		MaskCache:          MaskCacheStr,
		MaskCircuitBreaker: MaskCircuitBreakerStr,
		MaskHealth:         MaskHealthStr,
		MaskSecurity:       MaskSecurityStr,
		MaskPerformance:    MaskPerformanceStr,
		MaskInternal:       MaskInternalStr,
		MaskPerUser:        MaskPerUserStr,
		MaskPerRequest:     MaskPerRequestStr,
		MaskDetailed:       MaskDetailedStr,
	}

	for flag, name := range flags {
		if m.Has(flag) {
			parts = append(parts, name)
		}
	}

	return strings.Join(parts, "|")
}
