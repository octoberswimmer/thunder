package components

import (
	"testing"
)

// TestVerticalProgressBasic verifies that VerticalProgress returns a non-nil component.
func TestVerticalProgressBasic(t *testing.T) {
	steps := []ProgressStep{
		{Name: "Step 1", IsActive: true, IsCompleted: false},
		{Name: "Step 2", IsActive: false, IsCompleted: false},
	}
	comp := VerticalProgress(steps)
	if comp == nil {
		t.Error("VerticalProgress returned nil for valid steps")
	}
}

// TestVerticalProgressEmpty verifies that VerticalProgress works with empty steps.
func TestVerticalProgressEmpty(t *testing.T) {
	comp := VerticalProgress([]ProgressStep{})
	if comp == nil {
		t.Error("VerticalProgress returned nil for empty steps")
	}
}

// TestVerticalProgressStepStates verifies that all step states are handled.
func TestVerticalProgressStepStates(t *testing.T) {
	steps := []ProgressStep{
		{Name: "Completed", IsActive: false, IsCompleted: true},
		{Name: "Active", IsActive: true, IsCompleted: false},
		{Name: "Inactive", IsActive: false, IsCompleted: false},
	}
	comp := VerticalProgress(steps)
	if comp == nil {
		t.Error("VerticalProgress returned nil for mixed step states")
	}
}
