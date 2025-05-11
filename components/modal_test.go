package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestModalBasic verifies that Modal returns a non-nil component for a title and body.
func TestModalBasic(t *testing.T) {
	comp := Modal("Test Modal", masc.Text("Body content"))
	if comp == nil {
		t.Error("Modal returned nil for valid title and body")
	}
}
