package components

import (
	"testing"
)

// TestStaticField verifies that StaticField returns non-nil for valid parameters.
func TestStaticField(t *testing.T) {
	comp := StaticField("Test Label", "Test Value")
	if comp == nil {
		t.Error("StaticField returned nil for valid parameters")
	}
}

// TestStaticFieldEmpty verifies that StaticField returns non-nil for empty parameters.
func TestStaticFieldEmpty(t *testing.T) {
	comp := StaticField("", "")
	if comp == nil {
		t.Error("StaticField returned nil for empty parameters")
	}
}
