package components

import (
	"testing"

	"github.com/octoberswimmer/masc"
)

// TestPromptBasic verifies that Prompt returns a non-nil component for a title, body and actions.
func TestPromptBasic(t *testing.T) {
	comp := Prompt("Test Prompt", []masc.MarkupOrChild{Button("Okay", VariantNeutral, func(*masc.Event) {})}, masc.Text("Body content"))
	if comp == nil {
		t.Error("Prompt returned nil for valid title, body and actons")
	}
}
