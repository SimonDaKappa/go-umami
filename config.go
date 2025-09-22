package umami

import (
	"encoding/json"
	"os"
	"strings"
)

// Config represents the metrics configuration
type Config struct {
	// Global settings
	GlobalLevel Level `json:"global_level" yaml:"global_level"`

	// Per-group settings
	Groups map[string]GroupConfig `json:"groups" yaml:"groups"`

	// Backend configuration
	Backend BackendConfig `json:"backend" yaml:"backend"`
}

// GroupConfig represents configuration for a specific metric group
type GroupConfig struct {
	// Minimum level for this group
	Level Level `json:"level" yaml:"level"`
	// Level options for this group
	LevelOpts LevelOpts `json:"level_opts" yaml:"level_opts"`
}

// BackendConfig represents backend-specific configuration
type BackendConfig struct {
	Name   string         `json:"name" yaml:"name"`     // "prometheus", "datadog", "opentelemetry", etc.
	Config map[string]any `json:"config" yaml:"config"` // Backend-specific configuration
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		GlobalLevel: LevelImportant,
		Groups:      make(map[string]GroupConfig),
		Backend: BackendConfig{
			Name:   BackendNoneName,
			Config: make(map[string]any),
		},
	}
}

// ProductionConfig returns a production-ready configuration
func ProductionConfig(backend Backend) *Config {
	config := DefaultConfig()
	config.GlobalLevel = LevelImportant

	// Disable debug metrics in production
	for name, group := range config.Groups {
		if group.Level > LevelImportant {
			group.Level = LevelImportant
		}
		config.Groups[name] = group
	}

	config.Backend.Name = backend.Name()

	return config
}

// DevelopmentConfig returns a development configuration with more verbose metrics
func DevelopmentConfig(backend Backend) *Config {
	config := DefaultConfig()
	config.GlobalLevel = LevelVerbose

	// Enable detailed metrics in development
	for name, group := range config.Groups {
		group.Level = LevelVerbose
		config.Groups[name] = group
	}

	config.Backend.Name = backend.Name()

	return config
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

const (
	EnvMetricsBackendKey  string = "METRICS_BACKEND"
	EnvMetricsLevelKey    string = "METRICS_LEVEL"
	EnvMetricsMaskKey     string = "METRICS_MASK"
	EnvMetricsGroupPrefix string = "METRICS_GROUP_"
)

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {

	config := DefaultConfig()

	// Global level
	if levelStr := os.Getenv(EnvMetricsLevelKey); levelStr != "" {
		config.GlobalLevel = ParseLevel(levelStr)
	}

	// Backend type
	if backendType := os.Getenv(EnvMetricsBackendKey); backendType != "" {
		config.Backend.Name = backendType
	}

	// Group-specific overrides
	// Format: METRICS_GROUP_<NAME>_LEVEL=<level>
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, EnvMetricsGroupPrefix) {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := parts[1]

			keyParts := strings.Split(key, "_")
			if len(keyParts) < 4 {
				continue
			}

			groupName := strings.ToLower(keyParts[2])
			setting := strings.ToLower(keyParts[3])

			groupConfig := config.Groups[groupName]

			switch setting {
			case "level":
				groupConfig.Level = ParseLevel(value)
			}

			config.Groups[groupName] = groupConfig
		}
	}

	return config
}

// ApplyConfig applies the configuration to a metrics [Registry]
func ApplyConfig(manager Registry, config *Config) {
	globalLevelOpts := LevelOpts{
		ReplaceNoops: false,
	}

	// Apply global settings
	manager.SetGlobalLevel(config.GlobalLevel, globalLevelOpts)

	// Apply group-specific settings
	for name, groupConfig := range config.Groups {
		group := manager.Group(name)
		group.SetGroupLevel(groupConfig.Level, groupConfig.LevelOpts)
	}
}

// SaveToFile saves configuration to a JSON file
func (c *Config) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
