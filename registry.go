package umami

//--------------------------------------------------------------------------------
// File: registry.go
//
// This file contains the definition and implementation of the [Registry]
// interface for the umami metrics library.
//
// The Registry is responsible for managing [Group]s, global settings,
// level management, and Backend integration.
//
// It provides a central point for creating and accessing metric groups,
// each of which can have its own context and factory for creating metrics.
//
// A registry is not tied to any specific backend, but rather uses a pluggable
// [Backend] interface for each [Group] to create metric instances.
//--------------------------------------------------------------------------------

import (
	"slices"
	"sync"
)

// Registry is a global level management interface for metrics
//
// It is responsible for managing [Group]s and global settings,
// level management, and Backend integration.
type Registry interface {
	// Group returns a [Group] if it exists, or nil if it does not
	Group(name string) Group

	// NewGroup creates a new metric [Group] with the given name, [Backend], and [Level].
	//
	// If a group with the same name already exists, it is returned instead.
	NewGroup(name string, backend Backend, level ...Level) Group

	// SetGlobalLevel sets the global metrics level
	SetGlobalLevel(level Level, opts LevelOpts)

	// GlobalContext returns the global metrics context
	GlobalContext() Context
}

// registry implements the [Registry] interface
type registry struct {
	mu          sync.RWMutex
	groups      map[string]*group // Map of group name to group
	globalLevel Level
}

// NewRegistry creates a new metrics registry with the specified [Backend]
func NewRegistry(level Level) Registry {
	return &registry{
		groups:      make(map[string]*group),
		globalLevel: level,
	}
}

// NewGroup creates a new metric [Group] with the given name, [Backend], and [Level].
// If a group with the same name already exists, it is returned instead.
//
// It optionally accepts a variable number of [Level] arguments to set the minimum
// level for the group. If no level is provided, the registry's global level is used.
// Of those provided, the lowest level is chosen as the group's level.
//
// Note: This means that if a different group with a same name but different
// backend or level is requested, the existing group is returned and the new
// parameters are ignored.
func (m *registry) NewGroup(name string, backend Backend, level ...Level) Group {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group, exists := m.groups[name]; exists {
		return group
	}

	if len(level) == 0 {
		level = []Level{m.globalLevel}
	}
	minLevel := slices.Min(level)

	group := newGroup(backend, name, minLevel)
	m.groups[name] = group
	return group
}

// Group returns a metric [Group] if it exists, or nil if it does not
func (m *registry) Group(name string) Group {
	m.mu.Lock()
	defer m.mu.Unlock()

	if group, exists := m.groups[name]; exists {
		return group
	}

	return nil
}

// SetGlobalLevel sets the global metrics level
func (m *registry) SetGlobalLevel(level Level, opts LevelOpts) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.globalLevel = level
	// Update all existing groups
	for _, group := range m.groups {
		group.SetGroupLevel(level, opts)
	}
}

// GlobalContext returns the global metrics context
func (m *registry) GlobalContext() Context {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return NewContext(m.globalLevel)
}
