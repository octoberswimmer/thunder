package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestTextInputBasic verifies that TextInput returns non-nil for valid params.
func TestTextInputBasic(t *testing.T) {
	comp := TextInput("Name", "", "Enter name", func(e *masc.Event) {})
	if comp == nil {
		t.Error("TextInput returned nil for valid parameters")
	}
}
