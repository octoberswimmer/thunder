package components

import (
	"testing"

	"github.com/octoberswimmer/masc"
)

// TestValidatedTimepickerBasic verifies that ValidatedTimepicker returns non-nil for valid params.
func TestValidatedTimepickerBasic(t *testing.T) {
	validation := ValidationState{Required: false}

	comp := ValidatedTimepicker("Test Time", "14:30", validation, func(*masc.Event) {})
	if comp == nil {
		t.Error("ValidatedTimepicker returned nil for valid parameters")
	}
}

// TestValidatedTimepickerWithError verifies that ValidatedTimepicker returns non-nil when validation has errors.
func TestValidatedTimepickerWithError(t *testing.T) {
	validation := ValidationState{
		Required:     true,
		HasError:     true,
		ErrorMessage: "Time is required",
	}

	comp := ValidatedTimepicker("Required Time", "", validation, func(*masc.Event) {})
	if comp == nil {
		t.Error("ValidatedTimepicker returned nil for parameters with validation errors")
	}
}

// TestValidatedTimepickerRequired verifies that ValidatedTimepicker handles required fields.
func TestValidatedTimepickerRequired(t *testing.T) {
	validation := ValidationState{
		Required: true,
		HelpText: "Please select a time",
	}

	comp := ValidatedTimepicker("Required Time", "", validation, func(*masc.Event) {})
	if comp == nil {
		t.Error("ValidatedTimepicker returned nil for required field")
	}
}

// TestValidatedTimepickerWithTooltip verifies that ValidatedTimepicker handles tooltip support.
func TestValidatedTimepickerWithTooltip(t *testing.T) {
	validation := ValidationState{
		Required: false,
		Tooltip:  "Select preferred appointment time",
		HelpText: "Time should be during business hours",
	}

	comp := ValidatedTimepicker("Appointment Time", "09:00", validation, func(*masc.Event) {})
	if comp == nil {
		t.Error("ValidatedTimepicker returned nil with tooltip")
	}
}

// TestValidatedTimepickerNilHandler verifies that ValidatedTimepicker handles nil onChange handler.
func TestValidatedTimepickerNilHandler(t *testing.T) {
	validation := ValidationState{Required: false}

	comp := ValidatedTimepicker("Test", "12:00", validation, nil)
	if comp == nil {
		t.Error("ValidatedTimepicker returned nil with nil onChange handler")
	}
}
