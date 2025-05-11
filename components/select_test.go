package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestSelectBasic verifies that Select returns a non-nil component for valid parameters.
func TestSelectBasic(t *testing.T) {
	opts := []SelectOption{
		{Label: "One", Value: "1"},
		{Label: "Two", Value: "2"},
	}
	comp := Select("Choose an option", opts, "1", func(e *masc.Event) {})
	if comp == nil {
		t.Error("Select returned nil for valid parameters")
	}
}
