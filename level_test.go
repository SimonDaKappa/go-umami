package umami

import (
	"testing"
)

func TestBasicLevelSystem(t *testing.T) {
	backend := NewMockBackend()
	manager := NewRegistry(LevelDebug)

	// Test basic manager functionality
	globalCtx := manager.GlobalContext()
	if !globalCtx.Enabled(LevelDebug) {
		t.Error("Manager should default to LevelDebug")
	}

	// Test group creation and inheritance
	group := manager.NewGroup("test_group", backend)
	groupCtx := group.Context()
	if !groupCtx.Enabled(LevelDebug) {
		t.Error("Group should inherit manager level")
	}

	// Test metric creation and basic operations
	counter := group.Counter(
		CounterOpts{
			BasicMetricOpts: BasicMetricOpts{
				FromComposite: false,
			},
			MetricInfo: MetricInfo{Name: "test", Help: "Test"},
		},
		LevelCritical,
	)

	if err := counter.Inc(groupCtx); err != nil {
		t.Errorf("Counter operation failed: %v", err)
	}

	t.Log("Basic level system test passed")
}

func TestFactoryLevelFiltering(t *testing.T) {
	backend := NewMockBackend()
	manager := NewRegistry(LevelDebug)
	group := manager.NewGroup("test_group", backend)

	// Set group to only allow Critical metrics
	group.SetGroupLevel(LevelCritical, LevelOpts{ReplaceNoops: false})

	// Try to create metrics at different levels
	criticalCounter := group.Counter(
		CounterOpts{
			MetricInfo:      MetricInfo{Name: "critical", Help: "Critical"},
			BasicMetricOpts: BasicMetricOpts{FromComposite: false},
		},
		LevelCritical,
	)
	debugCounter := group.Counter(
		CounterOpts{
			MetricInfo:      MetricInfo{Name: "debug", Help: "Debug"},
			BasicMetricOpts: BasicMetricOpts{FromComposite: false},
		},
		LevelDebug,
	)

	groupCtx := group.Context()

	// Both should work at the interface level (no errors)
	if err := criticalCounter.Inc(groupCtx); err != nil {
		t.Errorf("Critical counter should work: %v", err)
	}
	if err := debugCounter.Inc(groupCtx); err != nil {
		t.Errorf("Debug counter should work (as noop): %v", err)
	}

	t.Log("Factory level filtering test passed")
}

func TestGroupLevelUpdates(t *testing.T) {
	backend := NewMockBackend()
	manager := NewRegistry(LevelDebug)
	group := manager.NewGroup("test_group", backend)

	// Create a metric
	counter := group.Counter(
		CounterOpts{
			MetricInfo:      MetricInfo{Name: "test_counter", Help: "Test"},
			BasicMetricOpts: BasicMetricOpts{FromComposite: false},
		},
		LevelDebug,
	)

	// Use it with group context
	groupCtx := group.Context()
	if err := counter.Inc(groupCtx); err != nil {
		t.Errorf("Counter should work with group context: %v", err)
	}

	// Change group level to be more restrictive
	group.SetGroupLevel(LevelCritical, LevelOpts{ReplaceNoops: false})

	// The counter should have its level updated - verify it works with critical context
	criticalCtx := NewContext(LevelCritical)
	if err := counter.Inc(criticalCtx); err != nil {
		t.Errorf("Counter should work with critical context after level update: %v", err)
	}

	t.Log("Group level updates test passed")
}
