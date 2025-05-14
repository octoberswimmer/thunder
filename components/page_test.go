package components

import (
	"github.com/octoberswimmer/masc"
	"testing"
)

// TestPageBasic verifies that Page returns a non-nil component for header and content.
func TestPageBasic(t *testing.T) {
	header := masc.Text("Header")
	content := masc.Text("Content")
	comp := Page(header, content)
	if comp == nil {
		t.Error("Page returned nil for valid header and content")
	}
}
