package umami

import (
	"testing"
)

// TestSwitchableMetrics demonstrates how the switchable wrapper pattern works
func TestSwitchableMetrics(t *testing.T) {
	// Create a mock backend
	backend := &mockBackend{}

	// Create a group with a level that disables metrics
	group := newGroup(backend, "test", LevelDisabled)

	// Create a counter - this should return a switchable wrapper containing a noop
	counter := group.Counter(
		CounterOpts{
			MetricInfo: MetricInfo{
				Help: "A test counter",
				Name: "test_counter",
			},
			BasicMetricOpts: BasicMetricOpts{FromComposite: false},
		},
		LevelDebug,
	)

	// Verify it's initially a noop by checking the switchable wrapper
	if switchable, ok := counter.(*switchableCounter); ok {
		if !switchable.IsNoop() {
			t.Error("Expected counter to initially be noop")
		}
	} else {
		t.Error("Expected counter to be switchable")
	}

	// The user keeps their reference to 'counter'
	// Now change the group level to enable metrics and convert noops
	group.SetGroupLevel(LevelDebug, LevelOpts{ReplaceNoops: true})

	// The same 'counter' reference should now be backed by a real implementation
	if switchable, ok := counter.(*switchableCounter); ok {
		if switchable.IsNoop() {
			t.Error("Expected counter to be converted to real implementation")
		}
	}

	// Test that the counter works
	ctx := NewContext(LevelDebug)
	err := counter.Inc(ctx)
	if err != nil {
		t.Errorf("Counter.Inc() failed: %v", err)
	}
}
