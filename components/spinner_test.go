package components

import "testing"

// TestSpinnerDefault verifies that Spinner returns a component for default size.
func TestSpinnerDefault(t *testing.T) {
	if comp := Spinner(""); comp == nil {
		t.Error("Spinner returned nil for default size")
	}
}

// TestSpinnerSizes verifies that Spinner handles different sizes.
func TestSpinnerSizes(t *testing.T) {
	for _, sz := range []string{"small", "medium", "large"} {
		if comp := Spinner(sz); comp == nil {
			t.Errorf("Spinner returned nil for size %s", sz)
		}
	}
}
