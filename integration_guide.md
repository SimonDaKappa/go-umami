# Integration Guide: New Metrics System with Existing Code

## Analysis of Current Implementation

After examining your middleware and AMQP pipeline code, here are the key integration opportunities:

## Current Problems

### 1. Middleware (`middleware.go`)
```go
// Current: Tightly coupled to concrete type
type WebMetricsHandler struct {
    Metrics *metrics.WebMetricGroup  // Concrete type
    next    http.Handler
}

// Problems:
// - Direct dependency on WebMetricGroup
// - No early return for disabled metrics
// - Error handling around metrics that could fail
// - Hard to test without real metrics backend
```

### 2. Pipeline (`pipeline.go` + `queue_funcs.go`)
```go
// Current: Complex instrumentation patterns
func (pipe *PipelineManager) InstrumentAMPQPublishFunc(...)

// Problems:
// - Very complex wrapper functions
// - Lots of error handling: "failed to get message queue publisher metric"
// - Boolean flags: performMetrics = true/false
// - Manual metric recording with error checking
```

## Proposed Integration Approach

### 1. Dependency Injection Pattern

Instead of injecting concrete metric groups, inject the manager and let components create their own metrics:

```go
// Current
func NewWebServer(metrics *metrics.WebMetricGroup) *WebServer

// Proposed  
func NewWebServer(metricsManager metrics.Manager) *WebServer {
    webGroup := metricsManager.Group("web")
    factory := webGroup.Factory()
    ctx := webGroup.Context()
    
    return &WebServer{
        requestCounter: factory.Counter("requests_total", 
            metrics.LevelCritical, metrics.MaskCounters, "route"),
        requestLatency: factory.Histogram("request_duration_seconds",
            metrics.LevelCritical, metrics.MaskLatency, defaultBuckets, "route"),
        metricsCtx: ctx,
    }
}
```

### 2. Simplified Middleware

```go
type WebMetricsHandler struct {
    next           http.Handler
    requestCounter metrics.Counter   // Interface, not concrete
    requestLatency metrics.Histogram // Interface, not concrete  
    metricsCtx     metrics.Context   // For early return
}

func (h *WebMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Early return if metrics disabled - zero cost
    if !h.metricsCtx.Enabled(metrics.LevelCritical) {
        h.next.ServeHTTP(w, r)
        return
    }
    
    // Record request - never fails
    h.requestCounter.Inc(h.metricsCtx)
    
    // Time the request using built-in timer
    rww := NewResponseWriterWrapper(w)
    h.requestLatency.Time(h.metricsCtx, func() error {
        h.next.ServeHTTP(rww, r)
        return nil
    })
    
    // No error handling needed - metrics never fail
}
```

### 3. Simplified Pipeline Metrics

Replace the complex instrumentation with simple, direct calls:

```go
type PipelineManager struct {
    // Replace complex metric group with simple interfaces
    publishCounter  metrics.Counter
    processTimer    metrics.Timer
    jobsInProgress  metrics.Gauge
    metricsCtx      metrics.Context
}

func (pipe *PipelineManager) PublishPromptJob(ctx context.Context, msg *api.PromptInMessage) (uuid.UUID, error) {
    // Simple, direct metrics - early return built-in
    pipe.publishCounter.Inc(pipe.metricsCtx)
    pipe.jobsInProgress.Inc(pipe.metricsCtx)
    
    // Actual business logic
    err := pipe.doPublish(msg)
    
    if err != nil {
        pipe.jobsInProgress.Dec(pipe.metricsCtx)
        return uuid.Nil, err
    }
    
    return msg.MessageID, nil
}

func (pipe *PipelineManager) modelOutQueueHandler(msg amqp.Delivery) (uuid.UUID, error) {
    // Automatic timing with early return
    var result uuid.UUID
    var err error
    
    pipe.processTimer.Time(pipe.metricsCtx, func() error {
        result, err = pipe.processModelOut(msg)
        return err
    })
    
    // Clean up
    pipe.jobsInProgress.Dec(pipe.metricsCtx)
    return result, err
}
```

## Key Benefits of This Approach

### 1. **Zero Coupling**
- Components depend only on interfaces
- No imports of concrete metric types
- Easy to mock for testing

### 2. **Zero Cost When Disabled**
- Early return pattern built into interfaces
- No performance impact when metrics disabled
- No complex boolean flags needed

### 3. **Zero Error Handling**
- Metrics never fail in new system
- Remove all `performMetrics` boolean logic
- Remove all metric error checking

### 4. **Simplified Code**
```go
// Before: Complex instrumentation
func (pipe *PipelineManager) InstrumentAMPQPublishFunc(ctx context.Context, queue string, handler AMPQPublishFunc) AMPQPublishFunc {
    performMetrics := true
    pubMetric, err := pipe.metric.GetMessageQueuePublisherMetric(queue)
    if err != nil {
        logging.Errorf("failed to get message queue publisher metric: %v", err)
        performMetrics = false
    }
    var pipeMetric *metrics.FullPipelineMetric = nil
    // ... 20+ lines of complex setup
}

// After: Simple direct calls  
func (pipe *PipelineManager) PublishMessage(msg Message) error {
    pipe.publishCounter.Inc(pipe.metricsCtx)  // That's it!
    return pipe.actualPublish(msg)
}
```

## Migration Strategy

### Phase 1: Update Constructors
- Change constructors to take `metrics.Manager` instead of concrete groups
- Create metrics using factory pattern inside constructors

### Phase 2: Replace Metric Fields
- Replace concrete metric group fields with interface fields
- Add `metricsCtx metrics.Context` field to each component

### Phase 3: Simplify Metric Calls
- Remove all the complex instrumentation functions
- Replace with direct calls to metric interfaces
- Remove all metric error handling

### Phase 4: Configuration
- Add level/mask configuration to your existing config system
- Set appropriate levels for production vs development

## Example Configuration Integration

```go
// In your existing config
type Config struct {
    // ... existing fields
    Metrics struct {
        Level string `json:"level"`
        Mask  string `json:"mask"`
        Groups map[string]struct {
            Level string `json:"level"`
            Mask  string `json:"mask"`
        } `json:"groups"`
    } `json:"metrics"`
}

// In main()
metricsConfig := metrics.LoadConfigurationFromEnv()
manager := metrics.NewManager(prometheusBackend)
metrics.ApplyConfiguration(manager, metricsConfig)

// Pass manager to components
webServer := web.NewServer(manager.Group("web"))
pipelineManager := pipeline.NewManager(manager.Group("pipeline"))
```

This approach eliminates the complex instrumentation patterns while providing better performance and maintainability.
