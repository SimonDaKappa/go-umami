package umami

import "sync"

// metricsContext implements the Context interface
type metricsContext struct {
	level Level
	mask  MetricMask
}

// NewContext creates a new metrics context with the given level and mask
func NewContext(level Level, mask MetricMask) Context {
	return &metricsContext{
		level: level,
		mask:  mask,
	}
}

// Enabled returns true if metrics at this level should be processed
func (c *metricsContext) Enabled(level Level) bool {
	return level.Enabled(c.level)
}

// EnabledMask returns true if metrics with this mask should be processed
func (c *metricsContext) EnabledMask(mask MetricMask) bool {
	return c.mask.Has(mask)
}

// WithLevel returns a new context with the specified level
func (c *metricsContext) WithLevel(level Level) Context {
	return &metricsContext{
		level: level,
		mask:  c.mask,
	}
}

// WithMask returns a new context with the specified mask
func (c *metricsContext) WithMask(mask MetricMask) Context {
	return &metricsContext{
		level: c.level,
		mask:  mask,
	}
}

// manager implements the Manager interface
type manager struct {
	mu          sync.RWMutex
	groups      map[string]*group
	globalLevel Level
	globalMask  MetricMask
	backend     Backend // Pluggable backend (prometheus, datadog, etc.)
}

// NewManager creates a new metrics manager with the specified backend
func NewManager(backend Backend) Manager {
	return &manager{
		groups:      make(map[string]*group),
		globalLevel: LevelImportant, // Safe default
		globalMask:  MaskProduction, // Safe default
		backend:     backend,
	}
}

// Group returns or creates a metric group
func (m *manager) Group(name string) Group {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group, exists := m.groups[name]; exists {
		return group
	}

	group := &group{
		name:    name,
		level:   m.globalLevel,
		mask:    m.globalMask,
		backend: m.backend,
	}
	m.groups[name] = group
	return group
}

// SetGlobalLevel sets the global metrics level
func (m *manager) SetGlobalLevel(level Level) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.globalLevel = level
	// Update all existing groups
	for _, group := range m.groups {
		group.SetLevel(level)
	}
}

// SetGlobalMask sets the global metrics mask
func (m *manager) SetGlobalMask(mask MetricMask) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.globalMask = mask
	// Update all existing groups
	for _, group := range m.groups {
		group.SetMask(mask)
	}
}

// GlobalContext returns the global metrics context
func (m *manager) GlobalContext() Context {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return NewContext(m.globalLevel, m.globalMask)
}

// group implements the Group interface
type group struct {
	mu      sync.RWMutex
	name    string
	level   Level
	mask    MetricMask
	backend Backend
	factory Factory
}

// Factory returns a factory for creating metrics in this group
func (g *group) Factory() Factory {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.factory == nil {
		g.factory = newFactory(g.backend, g.name, g.level, g.mask)
	}
	return g.factory
}

// SetLevel sets the level for this group
func (g *group) SetLevel(level Level) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.level = level
	// Invalidate factory so it gets recreated with new level
	g.factory = nil
}

// SetMask sets the mask for this group
func (g *group) SetMask(mask MetricMask) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.mask = mask
	// Invalidate factory so it gets recreated with new mask
	g.factory = nil
}

// Context returns a context for this group
func (g *group) Context() Context {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return NewContext(g.level, g.mask)
}
