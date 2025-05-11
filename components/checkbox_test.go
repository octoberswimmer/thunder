package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestCheckboxBasic verifies that Checkbox returns non-nil for valid parameters.
func TestCheckboxBasic(t *testing.T) {
	comp := Checkbox("Accept Terms", false, func(e *masc.Event) {})
	if comp == nil {
		t.Error("Checkbox returned nil for valid parameters")
	}
}
