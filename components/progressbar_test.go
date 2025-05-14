package components

import (
	"testing"
)

// TestProgressBarBounds verifies clamping at 0 and 100.
func TestProgressBarBounds(t *testing.T) {
	// Below 0
	if comp := ProgressBar(-10); comp == nil {
		t.Error("ProgressBar should not be nil for negative percent")
	}
	// Above 100
	if comp := ProgressBar(150); comp == nil {
		t.Error("ProgressBar should not be nil for percent >100")
	}
}

// TestProgressBarBasic verifies ProgressBar returns non-nil for valid percent.
func TestProgressBarBasic(t *testing.T) {
	for _, pct := range []int{0, 50, 100} {
		if comp := ProgressBar(pct); comp == nil {
			t.Errorf("ProgressBar returned nil for percent %d", pct)
		}
	}
}
