package components

import (
	"testing"
)

// TestValidatedRadioButtonGroupBasic verifies that ValidatedRadioButtonGroup returns non-nil for valid params.
func TestValidatedRadioButtonGroupBasic(t *testing.T) {
	options := []RadioButtonOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}
	validation := ValidationState{Required: false}

	comp := ValidatedRadioButtonGroup("Test Label", "test-group", options, "opt1", validation, func(string) {})
	if comp == nil {
		t.Error("ValidatedRadioButtonGroup returned nil for valid parameters")
	}
}

// TestValidatedRadioButtonGroupWithError verifies that ValidatedRadioButtonGroup returns non-nil when validation has errors.
func TestValidatedRadioButtonGroupWithError(t *testing.T) {
	options := []RadioButtonOption{
		{Label: "Option 1", Value: "opt1"},
	}
	validation := ValidationState{
		Required:     true,
		HasError:     true,
		ErrorMessage: "Test error",
	}

	comp := ValidatedRadioButtonGroup("Test Label", "test-group", options, "", validation, func(string) {})
	if comp == nil {
		t.Error("ValidatedRadioButtonGroup returned nil for parameters with validation errors")
	}
}

// TestValidatedRadioButtonGroupEmptyOptions verifies that ValidatedRadioButtonGroup handles empty options.
func TestValidatedRadioButtonGroupEmptyOptions(t *testing.T) {
	validation := ValidationState{Required: false}

	comp := ValidatedRadioButtonGroup("", "empty-group", []RadioButtonOption{}, "", validation, func(string) {})
	if comp == nil {
		t.Error("ValidatedRadioButtonGroup returned nil for empty options")
	}
}

// TestValidatedRadioButtonGroupRequired verifies that ValidatedRadioButtonGroup handles required fields.
func TestValidatedRadioButtonGroupRequired(t *testing.T) {
	options := []RadioButtonOption{
		{Label: "Required Option", Value: "req"},
	}
	validation := ValidationState{
		Required: true,
		HelpText: "This field is required",
	}

	comp := ValidatedRadioButtonGroup("Required Field", "req-group", options, "", validation, func(string) {})
	if comp == nil {
		t.Error("ValidatedRadioButtonGroup returned nil for required field")
	}
}

// TestValidatedRadioButtonGroupNilHandler verifies that ValidatedRadioButtonGroup handles nil onChange handler.
func TestValidatedRadioButtonGroupNilHandler(t *testing.T) {
	options := []RadioButtonOption{
		{Label: "Test", Value: "test"},
	}
	validation := ValidationState{Required: false}

	comp := ValidatedRadioButtonGroup("Test", "test-group", options, "test", validation, nil)
	if comp == nil {
		t.Error("ValidatedRadioButtonGroup returned nil with nil onChange handler")
	}
}
