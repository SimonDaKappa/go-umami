# Redesigned Metrics System

This is a complete redesign of the metrics system with the following goals:
- **Reduce coupling** between metrics and application code
- **Level-based enablement** similar to logging frameworks (slog/zap)
- **Early return pattern** for zero-cost disabled metrics
- **Pluggable backends** for different metrics libraries

## Key Features

### 1. Level-Based Metrics
Similar to log levels, metrics have importance levels:
- `LevelCritical` - Essential business metrics (always enabled)
- `LevelImportant` - Important operational metrics 
- `LevelDebug` - Debugging metrics (often disabled in production)
- `LevelVerbose` - Very detailed metrics (troubleshooting only)

### 2. Metric Masks
Fine-grained control over metric types:
- `MaskCounters`, `MaskLatency`, `MaskErrors` - Core metrics
- `MaskResources`, `MaskQueues`, `MaskConnections` - Infrastructure
- `MaskPerUser`, `MaskPerRequest` - High-cardinality metrics

### 3. Early Return Pattern
When metrics are disabled, calls return immediately with zero cost:

```go
func (c *counter) Inc(ctx Context) error {
    if !ctx.Enabled(c.level) || !ctx.EnabledMask(c.mask) {
        return nil // Early return - zero cost
    }
    return c.backend.Inc()
}
```

### 4. Interface-First Design
Application code only depends on interfaces, not concrete implementations:

```go
// Application code
func handleRequest(counter Counter, ctx Context) {
    counter.Inc(ctx) // Works with any backend
}
```

## Usage Examples

### Basic Setup
```go
// 1. Create backend (Prometheus, DataDog, etc.)
backend := NewPrometheusBackend()
manager := NewManager(backend)

// 2. Configure levels and masks
manager.SetGlobalLevel(LevelImportant)
manager.SetGlobalMask(MaskProduction)

// 3. Create metric groups
webGroup := manager.Group("web")
factory := webGroup.Factory()

// 4. Create metrics
counter := factory.Counter("requests_total", LevelCritical, MaskCounters, "method")
latency := factory.Histogram("request_duration", LevelCritical, MaskLatency, 
    []float64{0.001, 0.01, 0.1, 1}, "method")

// 5. Use metrics
ctx := webGroup.Context()
counter.Inc(ctx)
latency.Observe(ctx, 0.025)
```

### Runtime Configuration
```go
// Disable debug metrics in production
webGroup.SetLevel(LevelImportant)

// Enable only essential metrics
webGroup.SetMask(MaskEssential)
```

### Specialized Metrics
```go
// Database connection pool
dbPool := factory.Pool("connections", LevelImportant, MaskConnections)
dbPool.SetActive(ctx, 10)
dbPool.Acquired(ctx)

// Cache metrics
cache := factory.Cache("user_cache", LevelImportant, MaskCache)
cache.Hit(ctx)
cache.Miss(ctx)

// Circuit breaker
cb := factory.CircuitBreaker("external_api", LevelImportant, MaskCircuitBreaker)
cb.SetState(ctx, "open")
cb.Failure(ctx)
```

## Configuration

### Environment Variables
```bash
export METRICS_LEVEL=IMPORTANT
export METRICS_MASK=PRODUCTION
export METRICS_BACKEND=prometheus

# Per-group settings
export METRICS_GROUP_WEB_LEVEL=CRITICAL
export METRICS_GROUP_DEBUG_MASK=ALL
```

### JSON Configuration
```json
{
  "global_level": "IMPORTANT",
  "global_mask": "PRODUCTION",
  "groups": {
    "web": {
      "level": "CRITICAL",
      "mask": "COUNTERS|LATENCY|ERRORS"
    },
    "database": {
      "level": "IMPORTANT", 
      "mask": "CONNECTIONS|LATENCY"
    }
  },
  "backend": {
    "type": "prometheus",
    "config": {
      "namespace": "pacrag"
    }
  }
}
```

### Predefined Configurations
```go
// Production
config := ProductionConfiguration()

// Development 
config := DevelopmentConfiguration()

// Apply to manager
ApplyConfiguration(manager, config)
```

## Benefits

### 1. Zero Coupling
Application code only knows about interfaces:
```go
type WebHandler struct {
    requestCounter Counter  // Interface, not concrete type
    requestLatency Histogram
}
```

### 2. Zero Cost When Disabled
Disabled metrics have no performance impact:
```go
// If disabled, this returns immediately
counter.Inc(ctx) // No allocations, no backend calls
```

### 3. Pluggable Backends
Easy to swap between Prometheus, DataDog, StatsD, etc.:
```go
// Switch backends without changing application code
backend := NewDataDogBackend()  // or NewPrometheusBackend()
manager := NewManager(backend)
```

### 4. Flexible Configuration
Runtime control over which metrics are active:
```go
// Production: only critical metrics
manager.SetGlobalLevel(LevelCritical)

// Debugging: enable verbose metrics for specific component
webGroup.SetLevel(LevelVerbose)
```

### 5. Composable Metrics
Complex metrics are built from simple ones:
```go
// Cache metrics = hit counter + miss counter + size gauge
cache := factory.Cache("name", level, mask)
```

## Migration from Old System

1. Replace direct Prometheus calls with interface calls
2. Add level and mask to metric creation
3. Pass context to metric operations
4. Configure levels/masks instead of individual enable/disable

This design provides the flexibility and performance benefits you're looking for while maintaining clean separation of concerns.
