package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestDatepickerBasic verifies that Datepicker returns non-nil for valid parameters.
func TestDatepickerBasic(t *testing.T) {
	comp := Datepicker("Date", "", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Datepicker returned nil for valid parameters")
	}
}
