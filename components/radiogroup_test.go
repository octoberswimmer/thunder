package components

import "testing"

// TestRadioGroupBasic verifies that RadioGroup returns non-nil for valid params.
func TestRadioGroupBasic(t *testing.T) {
	opts := []RadioOption{
		{Label: "Contains", Value: "contains"},
		{Label: "Starts With", Value: "startswith"},
	}
	comp := RadioGroup("mode", "Filter Mode", opts, "contains", func(val string) {})
	if comp == nil {
		t.Error("RadioGroup returned nil for valid parameters")
	}
}
