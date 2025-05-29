package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestTimepickerBasic verifies that Timepicker returns non-nil for valid parameters.
func TestTimepickerBasic(t *testing.T) {
	comp := Timepicker("Appointment Time", "14:30", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Timepicker returned nil for valid parameters")
	}
}

// TestTimepickerEmpty verifies that Timepicker handles empty values and nil handlers.
func TestTimepickerEmpty(t *testing.T) {
	comp := Timepicker("Time", "", nil)
	if comp == nil {
		t.Error("Timepicker returned nil for empty value and nil handler")
	}
}
