package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestCardBasic verifies that Card returns a non-nil component for a simple title and body.
func TestCardBasic(t *testing.T) {
	comp := Card("Title", masc.Text("Body"))
	if comp == nil {
		t.Error("Card returned nil for valid title and body")
	}
}
