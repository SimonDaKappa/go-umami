package umami

//--------------------------------------------------------------------------------
// File: context.go
//
// This file contains the definition and implementation of the [Context] interface
// for the umami metrics library.
//--------------------------------------------------------------------------------

//--------------------------------------------------------------------------------
// Interfaces
//--------------------------------------------------------------------------------

// Context allows checking if metrics are enabled without coupling to specific
// implementations
//
// Each [Group] typically has its own Context, but they may be shared as a global
// context for the [Registry]
type Context interface {
	// Enabled returns true if metrics at this level should be processed
	Enabled(level Level) bool

	// WithLevel returns a new context curried with the specified level
	WithLevel(level Level) Context
}

//--------------------------------------------------------------------------------
// Context Implementation
//--------------------------------------------------------------------------------

// metricsContext implements the [Context] interface
type metricsContext struct {
	level Level
}

// NewContext creates a new [metricsContext] with the given [Level]
func NewContext(level Level) Context {
	return &metricsContext{
		level: level,
	}
}

// Enabled returns true if metrics at this level should be processed
func (c *metricsContext) Enabled(level Level) bool {
	return level.Enabled(c.level)
}

// WithLevel returns a new context with the specified level
func (c *metricsContext) WithLevel(level Level) Context {
	return &metricsContext{
		level: level,
	}
}
