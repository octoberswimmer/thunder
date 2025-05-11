package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestToastBasic verifies that Toast returns a non-nil component for valid inputs.
func TestToastBasic(t *testing.T) {
	comp := Toast(VariantSuccess, "Success!", "Operation completed.", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Toast returned nil for valid parameters")
	}
}
