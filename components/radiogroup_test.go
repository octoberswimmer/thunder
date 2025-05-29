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

// TestRadioButtonGroupBasic verifies that RadioButtonGroup returns non-nil for valid params.
func TestRadioButtonGroupBasic(t *testing.T) {
	opts := []RadioButtonOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
		{Label: "Option 3", Value: "opt3"},
	}
	comp := RadioButtonGroup("choices", opts, "opt1", func(val string) {})
	if comp == nil {
		t.Error("RadioButtonGroup returned nil for valid parameters")
	}
}

// TestRadioButtonGroupEmptyOptions verifies that RadioButtonGroup handles empty options.
func TestRadioButtonGroupEmptyOptions(t *testing.T) {
	comp := RadioButtonGroup("empty", []RadioButtonOption{}, "", nil)
	if comp == nil {
		t.Error("RadioButtonGroup returned nil for empty options")
	}
}

// TestRadioButtonGroupNilHandler verifies that RadioButtonGroup handles nil onChange handler.
func TestRadioButtonGroupNilHandler(t *testing.T) {
	opts := []RadioButtonOption{
		{Label: "Test", Value: "test"},
	}
	comp := RadioButtonGroup("test", opts, "test", nil)
	if comp == nil {
		t.Error("RadioButtonGroup returned nil with nil onChange handler")
	}
}
