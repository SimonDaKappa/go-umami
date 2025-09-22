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
//
// A metric with level L should be processed if
// L <= configuredLevel and configuredLevel != LevelDisabled
func (l Level) Enabled(configuredLevel Level) bool {
	return l <= configuredLevel && configuredLevel != LevelDisabled
}

type LevelOpts struct {
	ReplaceNoops bool // If true, replace no-op metrics when changing level
}
