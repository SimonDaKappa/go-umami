package umami

import (
	"encoding/json"
	"os"
	"strings"
)

// Configuration represents the metrics configuration
type Configuration struct {
	// Global settings
	GlobalLevel Level      `json:"global_level" yaml:"global_level"`
	GlobalMask  MetricMask `json:"global_mask" yaml:"global_mask"`

	// Per-group settings
	Groups map[string]GroupConfig `json:"groups" yaml:"groups"`

	// Backend configuration
	Backend BackendConfig `json:"backend" yaml:"backend"`
}

// GroupConfig represents configuration for a specific metric group
type GroupConfig struct {
	Level Level      `json:"level" yaml:"level"`
	Mask  MetricMask `json:"mask" yaml:"mask"`
}

// BackendConfig represents backend-specific configuration
type BackendConfig struct {
	Type   string         `json:"type" yaml:"type"`     // "prometheus", "datadog", "statsd", etc.
	Config map[string]any `json:"config" yaml:"config"` // Backend-specific configuration
}

// DefaultConfiguration returns a sensible default configuration
func DefaultConfiguration(backend string) *Configuration {
	return &Configuration{
		GlobalLevel: LevelImportant,
		GlobalMask:  MaskProduction,
		Groups:      map[string]GroupConfig{},
		Backend: BackendConfig{
			Type: "prometheus",
			Config: map[string]any{
				"namespace": "pacrag",
				"subsystem": "webserver",
			},
		},
	}
}

// ProductionConfiguration returns a production-ready configuration
func ProductionConfiguration() *Configuration {
	config := DefaultConfiguration()
	config.GlobalLevel = LevelImportant
	config.GlobalMask = MaskProduction

	// Disable debug metrics in production
	for name, group := range config.Groups {
		if group.Level > LevelImportant {
			group.Level = LevelImportant
		}
		config.Groups[name] = group
	}

	return config
}

// DevelopmentConfiguration returns a development configuration with more verbose metrics
func DevelopmentConfiguration() *Configuration {
	config := DefaultConfiguration()
	config.GlobalLevel = LevelVerbose
	config.GlobalMask = MaskAll

	// Enable detailed metrics in development
	for name, group := range config.Groups {
		group.Level = LevelVerbose
		group.Mask = MaskAll
		config.Groups[name] = group
	}

	return config
}

// LoadConfigurationFromFile loads configuration from a JSON file
func LoadConfigurationFromFile(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadConfigurationFromEnv loads configuration from environment variables
func LoadConfigurationFromEnv() *Configuration {
	config := DefaultConfiguration()

	// Global level
	if levelStr := os.Getenv("METRICS_LEVEL"); levelStr != "" {
		config.GlobalLevel = ParseLevel(levelStr)
	}

	// Global mask
	if maskStr := os.Getenv("METRICS_MASK"); maskStr != "" {
		config.GlobalMask = ParseMask(maskStr)
	}

	// Backend type
	if backendType := os.Getenv("METRICS_BACKEND"); backendType != "" {
		config.Backend.Type = backendType
	}

	// Group-specific overrides
	// Format: METRICS_GROUP_<NAME>_LEVEL=<level>
	// Format: METRICS_GROUP_<NAME>_MASK=<mask>
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "METRICS_GROUP_") {
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
			case "mask":
				groupConfig.Mask = ParseMask(value)
			}

			config.Groups[groupName] = groupConfig
		}
	}

	return config
}

// ParseMask parses a mask string into a MetricMask
func ParseMask(s string) MetricMask {
	if s == "" {
		return MaskProduction
	}

	switch strings.ToUpper(s) {
	case "NONE":
		return MaskNone
	case "ESSENTIAL":
		return MaskEssential
	case "PRODUCTION":
		return MaskProduction
	case "ALL":
		return MaskAll
	}

	// Parse individual flags separated by |
	var mask MetricMask
	flags := strings.Split(strings.ToUpper(s), "|")

	for _, flag := range flags {
		flag = strings.TrimSpace(flag)
		switch flag {
		case "COUNTERS":
			mask = mask.Add(MaskCounters)
		case "LATENCY":
			mask = mask.Add(MaskLatency)
		case "THROUGHPUT":
			mask = mask.Add(MaskThroughput)
		case "ERRORS":
			mask = mask.Add(MaskErrors)
		case "RESOURCES":
			mask = mask.Add(MaskResources)
		case "QUEUES":
			mask = mask.Add(MaskQueues)
		case "CONNECTIONS":
			mask = mask.Add(MaskConnections)
		case "CACHE":
			mask = mask.Add(MaskCache)
		case "CIRCUIT_BREAKER":
			mask = mask.Add(MaskCircuitBreaker)
		case "HEALTH":
			mask = mask.Add(MaskHealth)
		case "SECURITY":
			mask = mask.Add(MaskSecurity)
		case "PERFORMANCE":
			mask = mask.Add(MaskPerformance)
		case "INTERNAL":
			mask = mask.Add(MaskInternal)
		case "PER_USER":
			mask = mask.Add(MaskPerUser)
		case "PER_REQUEST":
			mask = mask.Add(MaskPerRequest)
		case "DETAILED":
			mask = mask.Add(MaskDetailed)
		}
	}

	return mask
}

// ApplyConfiguration applies the configuration to a metrics manager
func ApplyConfiguration(manager Manager, config *Configuration) {
	// Apply global settings
	manager.SetGlobalLevel(config.GlobalLevel)
	manager.SetGlobalMask(config.GlobalMask)

	// Apply group-specific settings
	for name, groupConfig := range config.Groups {
		group := manager.Group(name)
		group.SetLevel(groupConfig.Level)
		group.SetMask(groupConfig.Mask)
	}
}

// SaveConfiguration saves configuration to a JSON file
func (c *Configuration) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
